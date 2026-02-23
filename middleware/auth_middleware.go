package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequireLogin ensures a valid session exists
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		role := session.Get("role")

		if user == nil || role == nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireRole restricts access by allowed roles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userRole := session.Get("role")

		if userRole == nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		// Unauthorized access → redirect home
		c.HTML(http.StatusForbidden, "home.html", gin.H{
			"error": "Access denied. You are not authorized for this section.",
		})
		c.Abort()
	}
}
