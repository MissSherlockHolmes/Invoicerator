package controllers

import (
	"invoicerator/config"
	"invoicerator/models"
	"log"
	"net/http"

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
	c.Redirect(http.StatusFound, "/profile")
}

func ShowHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func ShowSignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
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
	// Clear the session cookie
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	// Redirect to the home page
	c.Redirect(http.StatusFound, "/")
}
