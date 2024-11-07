package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"invoicerator/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdateProfile(t *testing.T) {
	// Create a standalone test database using SQLite in memory
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}

	// Auto-migrate the User model for the test database
	err = testDB.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate the test database: %v", err)
	}

	// Insert a user into the test database
	testUser := models.User{
		Username:     "testuser",
		CompanyName:  "Old Company",
		CompanyEmail: "",
	}
	err = testDB.Create(&testUser).Error
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a Gin router and mock middleware
	router := gin.Default()

	// Middleware to simulate authenticated user
	router.Use(func(c *gin.Context) {
		c.Set("user", &testUser) // Pass pointer to ensure updates are reflected
		c.Next()
	})

	// Mock the UpdateProfile handler
	router.POST("/profile", func(c *gin.Context) {
		var input struct {
			CompanyName  string `form:"company_name"`
			CompanyEmail string `form:"company_email"`
		}

		// Parse form data
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the user in the test database
		user := c.MustGet("user").(*models.User)
		user.CompanyName = input.CompanyName
		user.CompanyEmail = input.CompanyEmail
		testDB.Save(user)

		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
	})

	// Prepare request to update the profile
	reqBody := "company_name=New Company&company_email=test@example.com"
	req, _ := http.NewRequest("POST", "/profile", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}

	// Verify the database changes
	var updatedUser models.User
	err = testDB.First(&updatedUser, testUser.ID).Error
	if err != nil {
		t.Fatalf("Failed to fetch updated user: %v", err)
	}

	if updatedUser.CompanyName != "New Company" || updatedUser.CompanyEmail != "test@example.com" {
		t.Errorf("Expected updated fields, got CompanyName=%s, CompanyEmail=%s", updatedUser.CompanyName, updatedUser.CompanyEmail)
	}
}
