package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ShowUploadPDFPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload_pdf.html", nil)
}

func UploadPDF(c *gin.Context) {
	file, err := c.FormFile("pdf")
	if err != nil {
		c.String(http.StatusBadRequest, "No file is received")
		return
	}

	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusInternalServerError, "Unable to save the file")
		return
	}

	// Use GPT API or other methods to extract data
	// For example purposes, assume extractedData is obtained
	extractedData := "Extracted data from PDF"

	// Remove the uploaded file
	os.Remove(filePath)

	// Redirect to create_invoice with extracted data
	c.HTML(http.StatusOK, "create_invoice.html", gin.H{"InvoiceData": extractedData})
}
