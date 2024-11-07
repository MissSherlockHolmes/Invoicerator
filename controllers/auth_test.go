package controllers

import (
	"bytes"
	"html/template"
	"invoicerator/config"
	"invoicerator/middleware"
	"invoicerator/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	// Load the .env.test file
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf("Failed to load .env.test: %v", err)
	}

	// Set up the database
	config.ConnectDatabase()

	// Run the tests
	code := m.Run()

	// Cleanup: Remove the test database
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath != "" {
		os.Remove(dbPath)
	}

	os.Exit(code)
}

func TestPerformLogin(t *testing.T) {
	// Setup
	config.ConnectDatabase()
	// Create a dummy user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	user := models.User{Username: "testuser", Password: string(hashedPassword)}
	config.DB.Create(&user)

	// Mock request and response
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", PerformLogin)

	formData := url.Values{}
	formData.Set("username", "testuser")
	formData.Set("password", "testpassword")

	req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusFound, w.Code)
	}
}

func TestPerformLoginInvalidUser(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString("username=invaliduser&password=wrongpassword"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, but got %d", w.Code)
	}
}

func TestPerformSignup(t *testing.T) {
	router := setupRouter()

	// Define a sample request body
	reqBody := "username=newuser&password=newpassword"
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}
}

// setupRouter initializes the Gin router with all routes and middleware
func setupRouter() *gin.Engine {
	// Create a new Gin engine
	router := gin.Default()

	// Set mode to TestMode
	gin.SetMode(gin.TestMode)

	router.SetHTMLTemplate(template.Must(template.New("").Parse(`
	 {{define "login.html"}}Login Page Mock{{end}}
	 {{define "signup.html"}}Signup Page Mock{{end}}
	 {{define "profile.html"}}Profile Page Mock{{end}}
 `)))

	// Use middleware
	router.Use(middleware.SetUserStatus())

	// Load routes
	router.POST("/login", PerformLogin)
	router.POST("/signup", PerformSignup)
	router.POST("/profile", UpdateProfile)

	// Wrap GenerateInvoicePDF in a handler function
	router.POST("/generate_invoice_pdf", func(c *gin.Context) {
		// Mock user data (replace with actual logic for fetching user)
		user := models.User{
			Username:       "testuser",
			CompanyName:    "Test Company",
			CompanyAddress: "123 Test Street",
			CompanyPhone:   "123-456-7890",
		}
		isPreview := true // Set this based on your test needs

		// Call GenerateInvoicePDF
		_, err := GenerateInvoicePDF(c, user, isPreview)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate invoice PDF"})
			return
		}
		c.JSON(200, gin.H{"success": true})
	})

	return router
}
