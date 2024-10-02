package main

import (
	"invoicerator/config"
	"invoicerator/controllers"
	"invoicerator/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Set trusted proxies to nil to avoid proxy warnings
	router.SetTrustedProxies(nil)

	config.ConnectDatabase()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Routes

	// Public routes
	router.GET("/", controllers.ShowHomePage)
	router.GET("/login", controllers.ShowLoginPage)
	router.POST("/login", controllers.PerformLogin)
	router.GET("/signup", controllers.ShowSignupPage)
	router.POST("/signup", controllers.PerformSignup)
	router.GET("/logout", controllers.Logout)

	// Protected routes (require authentication)
	authorized := router.Group("/", middleware.AuthRequired)
	{
		authorized.GET("/profile", controllers.ShowProfilePage)
		authorized.POST("/profile", controllers.UpdateProfile)
		authorized.GET("/create_invoice", controllers.ShowCreateInvoicePage)
		authorized.POST("/create_invoice", controllers.CreateInvoice)
		authorized.GET("/upload_pdf", controllers.ShowUploadPDFPage)
		authorized.POST("/upload_pdf", controllers.UploadPDF)
	}

	// Saved invoices route (public, uses local cache)
	router.GET("/saved_invoices", controllers.ShowSavedInvoicesPage)

	// Use the PORT environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	router.Run(":" + port)
}
