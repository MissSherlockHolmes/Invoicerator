package main

import (
	"invoicerator/config"
	"invoicerator/controllers"
	"invoicerator/middleware"
	"invoicerator/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Set trusted proxies to nil
	router.SetTrustedProxies(nil)

	// Use the middleware from the middleware package
	router.Use(middleware.SetUserStatus())

	config.ConnectDatabase()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Public routes
	router.GET("/", controllers.ShowHomePage)
	router.GET("/login", controllers.ShowLoginPage)
	router.POST("/login", controllers.PerformLogin)
	router.GET("/signup", controllers.ShowSignupPage)
	router.POST("/signup", controllers.PerformSignup)
	router.GET("/logout", controllers.Logout)

	// Protected routes
	authorized := router.Group("/", middleware.AuthRequired)
	{
		authorized.GET("/profile", controllers.ShowProfilePage)
		authorized.POST("/profile", controllers.UpdateProfile)
		authorized.GET("/create_invoice", controllers.ShowCreateInvoicePage)
		authorized.POST("/create_invoice", controllers.CreateInvoice)
		authorized.GET("/upload_pdf", controllers.ShowUploadPDFPage)
		authorized.POST("/upload_pdf", controllers.UploadPDF)
		authorized.GET("/options", controllers.ShowOptionsPage)
	}

	// Handle invoice preview generation
	router.POST("/preview_invoice", func(c *gin.Context) {
		// Retrieve user info from session
		username, err := c.Cookie("session_token")
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		var user models.User
		if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.String(http.StatusInternalServerError, "User not found")
			return
		}

		// Generate the PDF with user info and form data
		pdfData, err := controllers.GenerateInvoicePDF(c, user, true) // Pass 'true' for preview mode
		if err != nil {
			c.String(http.StatusInternalServerError, "Error generating PDF")
			return
		}

		// Return the PDF for preview
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "inline; filename=invoice_preview.pdf")
		c.Writer.Write(pdfData)
	})

	// Handle final invoice creation and sending via SendGrid
	router.POST("/send_invoice", func(c *gin.Context) {
		// Retrieve user info from session
		username, err := c.Cookie("session_token")
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		var user models.User
		if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.String(http.StatusInternalServerError, "User not found")
			return
		}

		// Generate the final PDF invoice
		pdfData, err := controllers.GenerateInvoicePDF(c, user, false)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error generating PDF")
			return
		}

		// Send the invoice via SendGrid
		recipientEmail := c.PostForm("email")     // Client email from form
		userEmails := []string{user.CompanyEmail} // Add the user's company email
		if err := controllers.SendInvoiceWithSendGrid(pdfData, recipientEmail, userEmails); err != nil {
			c.String(http.StatusInternalServerError, "Error sending invoice via email")
			return
		}

		// Return success message
		c.String(http.StatusOK, "Invoice created and sent successfully")
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
