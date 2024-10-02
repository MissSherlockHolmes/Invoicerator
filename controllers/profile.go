package controllers

import (
	"encoding/json"
	"fmt"
	"invoicerator/config"
	"invoicerator/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowProfilePage(c *gin.Context) {
	username, _ := c.Cookie("session_token")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"Error": "User not found"})
		return
	}

	// Prepare the selected fields map
	selectedFieldsMap := make(map[string]bool)
	if user.SelectedFields != "" {
		var selectedFields []string
		json.Unmarshal([]byte(user.SelectedFields), &selectedFields)
		for _, field := range selectedFields {
			selectedFieldsMap[field] = true
		}
	}

	// Define available fields
	availableFields := []string{"Item Description", "Quantity", "Price", "Tax", "Discount"}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"User":              user,
		"AvailableFields":   availableFields,
		"SelectedFieldsMap": selectedFieldsMap,
	})
}

func UpdateProfile(c *gin.Context) {
	username, _ := c.Cookie("session_token")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"Error": "User not found"})
		return
	}

	// Get form data
	companyName := c.PostForm("company_name")
	companyAddress := c.PostForm("company_address")
	companyPhone := c.PostForm("company_phone")
	selectedFields := c.PostFormArray("fields")

	// Update user fields
	user.CompanyName = companyName
	user.CompanyAddress = companyAddress
	user.CompanyPhone = companyPhone

	// Save selected fields as JSON
	selectedFieldsJSON, _ := json.Marshal(selectedFields)
	user.SelectedFields = string(selectedFieldsJSON)

	// Handle letterhead upload
	letterheadFile, err := c.FormFile("letterhead")
	if err == nil {
		// Save the file
		letterheadPath := fmt.Sprintf("static/uploads/%s_%s", username, letterheadFile.Filename)
		err = c.SaveUploadedFile(letterheadFile, letterheadPath)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"Error": "Failed to upload letterhead"})
			return
		}
		user.LetterheadPath = "/" + letterheadPath // Adjust path for static serving
	}

	// Save user to database
	if err := config.DB.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"Error": "Failed to update profile"})
		return
	}

	// Rebuild selectedFieldsMap based on updated user data
	selectedFieldsMap := make(map[string]bool)
	if user.SelectedFields != "" {
		var updatedSelectedFields []string
		json.Unmarshal([]byte(user.SelectedFields), &updatedSelectedFields)
		for _, field := range updatedSelectedFields {
			selectedFieldsMap[field] = true
		}
	}

	// Define available fields
	availableFields := []string{"Item Description", "Quantity", "Price", "Tax", "Discount"}

	// Reload the profile page with success message
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"User":              user,
		"Success":           "Profile updated successfully",
		"AvailableFields":   availableFields,
		"SelectedFieldsMap": selectedFieldsMap,
	})
}
