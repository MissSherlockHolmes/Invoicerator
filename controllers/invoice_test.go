package controllers

import (
	"invoicerator/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestGenerateInvoicePDF tests the PDF generation function
func TestGenerateInvoicePDF(t *testing.T) {
	// Create a mock Gin context
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Add mock form data to the context
	c.Request = &http.Request{
		Method: "POST",
		Header: make(http.Header),
		PostForm: url.Values{
			"email":              {"test@example.com"},
			"recipient_name":     {"Test User"},
			"invoice_data":       {"Invoice #123"},
			"item_description[]": {"Test Item"},
			"item_quantity[]":    {"1.00"},
			"item_rate[]":        {"10.00"},
			"notes":              {"Thank you for your business."},
			"total_due":          {"10.00"},
		},
	}
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Mock user data
	user := models.User{
		Username:       "testuser",
		CompanyName:    "Test Company",
		CompanyAddress: "123 Test Street",
		CompanyPhone:   "123-456-7890",
	}

	// Call the function
	_, err := GenerateInvoicePDF(c, user, true)
	if err != nil {
		t.Errorf("Failed to generate invoice PDF: %v", err)
	}
}
