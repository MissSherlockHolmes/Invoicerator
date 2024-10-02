package controllers

import (
    "bytes"
    "encoding/base64"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jung-kurt/gofpdf"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

func ShowCreateInvoicePage(c *gin.Context) {
    c.HTML(http.StatusOK, "create_invoice.html", nil)
}

func CreateInvoice(c *gin.Context) {
    recipientEmail := c.PostForm("email")
    invoiceData := c.PostForm("invoice_data")

    // Generate PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Invoice")
    pdf.Ln(12)
    pdf.SetFont("Arial", "", 12)
    pdf.MultiCell(0, 10, invoiceData, "", "", false)

    // Optionally include company letterhead if available
    // pdf.ImageOptions(...)

    // Save PDF to buffer
    var buf bytes.Buffer
    err := pdf.Output(&buf)
    if err != nil {
        c.String(http.StatusInternalServerError, "Could not generate PDF")
        return
    }

    // Send email with SendGrid
    from := mail.NewEmail("Your Company", "no-reply@yourdomain.com")
    subject := "Your Invoice"
    to := mail.NewEmail("", recipientEmail)
    plainTextContent := "Please find your invoice attached."
    message := mail.NewSingleEmail(from, subject, to, plainTextContent, "")

    attachment := mail.NewAttachment()
    encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
    attachment.SetContent(encoded)
    attachment.SetType("application/pdf")
    attachment.SetFilename("invoice.pdf")
    attachment.SetDisposition("attachment")
    message.AddAttachment(attachment)

    client := sendgrid.NewSendClient("YOUR_SENDGRID_API_KEY")
    response, err := client.Send(message)
    if err != nil || response.StatusCode >= 400 {
        c.String(http.StatusInternalServerError, "Could not send email")
        return
    }

    // Store invoice in local cache (client-side)
    // Send data to client to store in localStorage
    c.Redirect(http.StatusFound, "/saved_invoices")
}
