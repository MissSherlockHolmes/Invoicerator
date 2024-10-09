package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"invoicerator/config"
	"invoicerator/models"
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
		pdf.ImageOptions("."+user.LetterheadPath, 10, 10, 0, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
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
		quantity, err := strconv.ParseFloat(itemQuantities[i], 64)
		if err != nil {
			return nil, err
		}

		rate, err := strconv.ParseFloat(itemRates[i], 64)
		if err != nil {
			return nil, err
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

	// Output the PDF to memory (using a buffer instead of a file)
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SendInvoiceWithSendGrid sends the invoice via SendGrid
func SendInvoiceWithSendGrid(pdfData []byte, recipientEmail string, userEmails []string) error {
	from := mail.NewEmail("Invoicerator", "no-reply@yourcompany.com")
	to := mail.NewEmail("Client", recipientEmail)

	subject := "Your Invoice"
	plainTextContent := "Please find your invoice attached."
	htmlContent := "<p>Thank you for your business.</p>"

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = subject

	personalization := mail.NewPersonalization()
	personalization.AddTos(to)

	// Add user emails as CC
	for _, email := range userEmails {
		cc := mail.NewEmail("", email)
		personalization.AddCCs(cc)
	}
	message.AddPersonalizations(personalization)

	// Set email content
	message.AddContent(mail.NewContent("text/plain", plainTextContent))
	message.AddContent(mail.NewContent("text/html", htmlContent))

	// Add PDF attachment
	encodedPDF := base64.StdEncoding.EncodeToString(pdfData)
	attachment := mail.NewAttachment()
	attachment.SetContent(encodedPDF)
	attachment.SetType("application/pdf")
	attachment.SetFilename("invoice.pdf")
	attachment.SetDisposition("attachment")
	message.AddAttachment(attachment)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		fmt.Println("Failed to send invoice:", err)
		return err
	}

	fmt.Printf("Email sent. Status code: %d\n", response.StatusCode)
	return nil
}

// CreateInvoice handles the final invoice creation and sends it via SendGrid
func CreateInvoice(c *gin.Context) {
	// Retrieve user info from session
	username, _ := c.Cookie("session_token")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, "User not found")
		return
	}

	// Generate the final PDF invoice in memory
	pdfData, err := GenerateInvoicePDF(c, user, false)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating PDF")
		return
	}

	// Get recipient email from form data
	recipientEmail := c.PostForm("email")     // Client email from form
	userEmails := []string{user.CompanyEmail} // Add the user's company email

	// Send the invoice via SendGrid
	if err := SendInvoiceWithSendGrid(pdfData, recipientEmail, userEmails); err != nil {
		c.String(http.StatusInternalServerError, "Error sending invoice via email")
		return
	}

	// Return success message
	c.String(http.StatusOK, "Invoice created and sent successfully")
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
