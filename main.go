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
		router.GET("/options", middleware.IsAuthenticated, controllers.ShowOptionsPage)

		// Add the options route here
		//authorized.GET("/options", controllers.ShowOptionsPage)
	}

	// Saved invoices route
	router.GET("/saved_invoices", controllers.ShowSavedInvoicesPage)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
