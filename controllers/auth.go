package controllers

import (
	"fmt"
	"invoicerator/config"
	"invoicerator/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Invalid username or password"})
		return
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Invalid username or password"})
		return
	}

	// Set session cookie
	c.SetCookie("session_token", username, 3600, "/", "", false, true)
	log.Println("User logged in successfully:", username)

	// Redirect to the profile page
	c.Redirect(http.StatusFound, "/options")
}

func ShowHomePage(c *gin.Context) {
	authenticated, _ := c.Get("Authenticated")
	username, _ := c.Get("Username")

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Authenticated": authenticated,
		"Username":      username,
	})
}

func ShowOptionsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "options.html", gin.H{})
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func ShowSignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func ShowCreateInvoicePage(c *gin.Context) {
	fmt.Println("Rendering Create Invoice Page") // Debugging
	c.HTML(http.StatusOK, "create_invoice.html", nil)
}

func PerformSignup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check if the username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"Error": "Username already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{"Error": "Error creating account"})
		return
	}

	// Create a new user
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	config.DB.Create(&user)

	// Set session cookie
	c.SetCookie("session_token", username, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/profile")
}

func Logout(c *gin.Context) {
	// Log the user out for debugging purposes
	sessionToken, err := c.Cookie("session_token")
	if err == nil {
		log.Println("Attempting to log out user:", sessionToken)
	}

	// Clear the session cookie (attempt multiple variations for Safari)
	domains := []string{"localhost", ""} // Try both "localhost" and "" (empty domain)
	paths := []string{"/", ""}           // Try both "/" and "" (root and empty paths)

	for _, domain := range domains {
		for _, path := range paths {
			// First clear with the default method
			c.SetCookie("session_token", "", -1, path, domain, false, true)

			// Then clear explicitly using http.SetCookie for more control
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     path,                  // Try both "/" and "" as path values
				Domain:   domain,                // Try both "localhost" and "" as domain values
				MaxAge:   -1,                    // Expire immediately
				Expires:  time.Unix(0, 0),       // Set expiration in the past
				Secure:   true,                  // Change to true if you're using HTTPS
				HttpOnly: true,                  // Ensure HttpOnly is true
				SameSite: http.SameSiteNoneMode, // Use None for Safari compatibility
			})
		}
	}

	// Log to verify if the cookie deletion was attempted
	log.Println("Logout attempt completed. Cookie should be cleared.")

	// Redirect to home page after clearing the cookie
	c.Redirect(http.StatusFound, "/")
}
