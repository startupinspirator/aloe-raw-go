package middleware

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login required"})
		c.Abort()
		return
	}
	c.Set("user_id", userID)
	c.Set("user_role", session.Get("user_role")) // Added role
	c.Next()
}

func RequireAdmin(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get("user_role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		c.Abort()
		return
	}
	c.Next()
}
