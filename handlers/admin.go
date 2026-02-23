package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ShowAdminDashboard renders the admin dashboard page
func ShowAdminDashboard(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get("role")
	if role != "admin" {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "admin_dashboard.html", nil)
}
