package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetUserStatus middleware sets the authentication status in the context
func SetUserStatus() gin.HandlerFunc {
    return func(c *gin.Context) {
        sessionToken, err := c.Cookie("session_token")
        if err != nil || sessionToken == "" {
            c.Set("Authenticated", false)
        } else {
            c.Set("Authenticated", true)
            c.Set("Username", sessionToken)
        }
        c.Next()
    }
}

// AuthRequired middleware checks if the user is authenticated
func AuthRequired(c *gin.Context) {
    sessionToken, err := c.Cookie("session_token")
    if err != nil || sessionToken == "" {
        c.Redirect(http.StatusFound, "/login")
        c.Abort()
        return
    }
    c.Next()
}

func IsAuthenticated(c *gin.Context) {
    _, err := c.Cookie("session_token")
    if err != nil {
        c.Redirect(http.StatusFound, "/login")
        return
    }
    c.Next()
}
