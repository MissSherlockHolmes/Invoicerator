package controllers

import (
    "fmt"
    "invoicerator/config"
    "invoicerator/models"
    "net/http"
    "github.com/gin-gonic/gin"
)

func ShowProfilePage(c *gin.Context) {
    username, err := c.Cookie("session_token")
    if err != nil {
        c.HTML(http.StatusUnauthorized, "profile.html", gin.H{"Error": "Not logged in"})
        return
    }

    var user models.User
    if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"Error": "User not found"})
        return
    }

    c.HTML(http.StatusOK, "profile.html", gin.H{
        "User": user,
    })
}

func UpdateProfile(c *gin.Context) {
    username, err := c.Cookie("session_token")
    if err != nil {
        c.HTML(http.StatusUnauthorized, "profile.html", gin.H{"Error": "Not logged in"})
        return
    }

    var user models.User
    if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"Error": "User not found"})
        return
    }

    // Get form data
    companyName := c.PostForm("company_name")
    companyEmail := c.PostForm("company_email")
    companyAddress := c.PostForm("company_address")
    companyPhone := c.PostForm("company_phone")
    termsConditions := c.PostForm("terms_conditions")

    // Update user fields
    user.CompanyName = companyName
    user.CompanyEmail = companyEmail
    user.CompanyAddress = companyAddress
    user.CompanyPhone = companyPhone
    user.TermsConditions = termsConditions

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

    // Reload the profile page with success message
    c.HTML(http.StatusOK, "profile.html", gin.H{
        "User":    user,
        "Success": "Profile updated successfully",
    })
}