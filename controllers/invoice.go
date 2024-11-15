package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"invoicerator/config"
	"invoicerator/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// GenerateInvoicePDF generates a PDF from form data and user profile data.
func GenerateInvoicePDF(c *gin.Context, user models.User, isPreview bool) ([]byte, error) {
	// Get form data
	recipientEmail := c.PostForm("email")
	recipientName := c.PostForm("recipient_name")
	invoiceData := c.PostForm("invoice_data")
	itemDescriptions := c.PostFormArray("item_description[]")
	itemQuantities := c.PostFormArray("item_quantity[]")
	itemRates := c.PostFormArray("item_rate[]")
	notes := c.PostForm("notes")
	totalDue := c.PostForm("total_due")

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add company letterhead/logo if available
	if user.LetterheadPath != "" {
		pdf.ImageOptions("."+user.LetterheadPath, 10, 10, 0, 20, false, gofpdf.ImageOptions{ImageType: "JPG"}, 0, "")
		pdf.Ln(20) // Add some space after the logo
	}

	// Add company information (from the user's profile)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, user.CompanyName)
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, user.CompanyAddress, "", "", false)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Phone: "+user.CompanyPhone)
	pdf.Ln(12)

	// Add recipient information
	pdf.Cell(40, 10, "Bill To: "+recipientName)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Email: "+recipientEmail)
	pdf.Ln(12)

	// Add invoice details
	pdf.MultiCell(0, 5, invoiceData, "", "", false)
	pdf.Ln(6)

	// Add line items
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(60, 10, "Description")
	pdf.Cell(20, 10, "Qty")
	pdf.Cell(40, 10, "Rate")
	pdf.Cell(40, 10, "Amount")
	pdf.Ln(10)

	for i := 0; i < len(itemDescriptions); i++ {
		quantityStr := itemQuantities[i]
		rateStr := itemRates[i]

		log.Printf("Processing item %d: quantity='%s', rate='%s'", i, quantityStr, rateStr)

		quantity, err := strconv.ParseFloat(quantityStr, 64)
		if err != nil {
			log.Printf("Error parsing quantity for item %d: %v", i, err)
			return nil, fmt.Errorf("invalid quantity for item %d: %w", i+1, err)
		}

		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			log.Printf("Error parsing rate for item %d: %v", i, err)
			return nil, fmt.Errorf("invalid rate for item %d: %w", i+1, err)
		}

		amount := quantity * rate
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(60, 10, itemDescriptions[i])
		pdf.Cell(20, 10, fmt.Sprintf("%.2f", quantity))
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", rate))
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", amount))
		pdf.Ln(10)
	}

	// Add notes section
	pdf.Ln(12)
	pdf.SetFont("Arial", "I", 10)
	pdf.MultiCell(0, 5, notes, "", "", false)

	// Add total due
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Total Due: $"+totalDue)


	// Add Terms and Conditions section if they exist
	if user.TermsConditions != "" {
		pdf.Ln(20)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 10, "Terms and Conditions")
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 9)
		pdf.MultiCell(0, 5, user.TermsConditions, "", "", false)
	}

	// Output the PDF to memory (using a buffer instead of a file)
	var buf bytes.Buffer
	log.Printf("generating pdf..")
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("ERROR: %s..", err)
		return nil, err
	}
	log.Printf("Done generating pdf.")

	return buf.Bytes(), nil
}

func SendInvoiceWithSendGrid(pdfData []byte, recipientName string, recipientEmail string, companyEmail string, companyName string) error {

	// Create the personalization object for the primary recipient
	personalization := mail.NewPersonalization()
	personalization.AddTos(mail.NewEmail(recipientName, recipientEmail)) // Primary recipient only

	message := mail.NewV3Mail()
	message.AddPersonalizations(personalization)

	if companyEmail != "" {
		cc := mail.NewEmail("", companyEmail)
		message.Personalizations[0].AddCCs(cc)
	}

	// Set the email sender and subject with company name
	from := mail.NewEmail(companyName, "invoices@invoicerator.com")
	subject := fmt.Sprintf("You have received an invoice from %s", companyName)
	message.SetFrom(from)
	message.Subject = subject

	// Create a more professional email content
	plainContent := fmt.Sprintf("You have received an invoice from %s. Please find it attached to this email.", companyName)
	htmlContent := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; color: #333;">
            <p>You have received an invoice from %s.</p>
            <p>The invoice is attached to this email as a PDF document.</p>
            <p>If you have any questions, please contact us directly.</p>
        </div>
    `, companyName)

	message.AddContent(mail.NewContent("text/plain", plainContent))
	message.AddContent(mail.NewContent("text/html", htmlContent))

	// Add the PDF attachment
	encodedPDF := base64.StdEncoding.EncodeToString(pdfData)
	attachment := mail.NewAttachment()
	attachment.SetContent(encodedPDF)
	attachment.SetType("application/pdf")
	attachment.SetFilename("invoice.pdf")
	attachment.SetDisposition("attachment")
	message.AddAttachment(attachment)

	// Send the email
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println("Failed to send invoice:", err)
		return err
	}
	log.Printf("Email sent. Status Code: %d, Response Body: %s", response.StatusCode, response.Body)
	return nil
}

func CreateInvoice(c *gin.Context) {
	// Retrieve user info from session
	username, _ := c.Cookie("session_token")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Generate the final PDF invoice in memory
	pdfData, err := GenerateInvoicePDF(c, user, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating PDF"})
		return
	}

	// Get recipient email from form data
	recipientEmail := c.PostForm("email") // Client email from form
	recipientName := c.PostForm("recipient_name")

	// Send the invoice via SendGrid
	if err := SendInvoiceWithSendGrid(pdfData, recipientName, recipientEmail, user.CompanyEmail, user.CompanyName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending invoice via email"})
		return
	}

	// Return success message as JSON
	c.JSON(http.StatusOK, gin.H{"message": "Invoice created and sent successfully"})
}

// PreviewInvoice handles the invoice preview generation
func PreviewInvoice(c *gin.Context) {
	// Retrieve user info from session
	username, _ := c.Cookie("session_token")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, "User not found")
		return
	}

	// Generate the preview PDF in memory
	pdfData, err := GenerateInvoicePDF(c, user, true)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating PDF preview")
		return
	}

	// Return the preview PDF as bytes
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=invoice_preview.pdf")
	c.Writer.Write(pdfData)
}
