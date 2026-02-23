package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionTimeout logs users out after 10 minutes (600 s) of inactivity.
func SessionTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		lastActivityRaw := session.Get("lastActivity")
		now := time.Now()

		if lastActivityRaw != nil {
			if lastActivity, ok := lastActivityRaw.(time.Time); ok {
				inactivity := now.Sub(lastActivity)
				if inactivity > 10*time.Minute {
					// Clear session → force re-login
					session.Clear()
					session.Save()
					c.HTML(http.StatusUnauthorized, "home.html", gin.H{
						"error": "Session expired due to 10 minutes of inactivity. Please log in again.",
					})
					c.Abort()
					return
				}
			}
		}

		// Update timestamp for active users
		session.Set("lastActivity", now)
		session.Save()

		c.Next()
	}
}
