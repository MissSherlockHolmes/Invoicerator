package controllers

import (
	"invoicerator/config"
	"invoicerator/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func ShowCreateInvoicePage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_invoice.html", nil)
}

func CreateInvoice(c *gin.Context) {
	// Existing code to get recipient email and invoice data

	username, _ := c.Cookie("session_token")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, "User not found")
		return
	}

	// Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add company letterhead if available
	if user.LetterheadPath != "" {
		pdf.ImageOptions("."+user.LetterheadPath, 10, 10, 0, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		pdf.Ln(20)
	}

	// Add company information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, user.CompanyName)
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, user.CompanyAddress, "", "", false)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Phone: "+user.CompanyPhone)
	pdf.Ln(12)

	// Rest of the invoice generation code

	// Existing code to send email
}
