package controllers

import (
	"invoicerator/config"
	"invoicerator/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowProfilePage(c *gin.Context) {
	username, _ := c.Cookie("session_token")
	var user models.User
	config.DB.Where("username = ?", username).First(&user)

	c.HTML(http.StatusOK, "profile.html", gin.H{"User": user})
}

func UpdateProfile(c *gin.Context) {
	username, _ := c.Cookie("session_token")
	var user models.User
	config.DB.Where("username = ?", username).First(&user)

	// Handle selected fields
	//selectedFields := c.PostFormArray("fields")
	// Convert selectedFields to a string or store as needed

	// Handle company letterhead upload
	letterheadFile, err := c.FormFile("letterhead")
	if err == nil {
		// Save the file temporarily or process as needed
		letterheadPath := "./static/uploads/" + letterheadFile.Filename
		c.SaveUploadedFile(letterheadFile, letterheadPath)
		// Update user record with letterhead path
	}

	// Update user in the database
	config.DB.Save(&user)

	c.Redirect(http.StatusFound, "/create_invoice")
}
