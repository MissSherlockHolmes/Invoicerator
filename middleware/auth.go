package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil || sessionToken == "" {
		log.Println("No valid session token found. Redirecting to login.")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	log.Println("Authenticated user with session token:", sessionToken)
	c.Next()
}
