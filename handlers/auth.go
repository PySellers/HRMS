package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"pysellers-erp-go/models"
	"pysellers-erp-go/security"
	"pysellers-erp-go/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var userFile = "data/db.json"

// Show login page with fresh captcha
// --- ShowHome: Generate and store captcha ---
func ShowHome(c *gin.Context) {
	session := sessions.Default(c)

	// Generate captcha
	captcha := utils.GenerateCaptcha()
	parts := strings.Split(captcha, "|")

	if len(parts) != 2 {
		c.String(http.StatusInternalServerError, "Captcha generation failed")
		return
	}

	// Save answer to session
	session.Set("captcha_answer", parts[0])
	if err := session.Save(); err != nil {
		fmt.Println("Session save failed:", err)
	}

	fmt.Println("Captcha stored:", parts[0]) // Debug confirmation

	c.HTML(http.StatusOK, "home.html", gin.H{
		"captchaQuestion": parts[1],
	})
}

// Secure Login Handler (with bcrypt, throttling, and fixed captcha)
func Login(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")
	role := strings.TrimSpace(c.PostForm("role"))
	captchaInput := strings.TrimSpace(c.PostForm("captcha"))

	session := sessions.Default(c)

	// --- Captcha Validation ---
	stored, _ := session.Get("captcha_answer").(string)
	// 🚨 DEMO MODE BYPASS — REMOVE AFTER DEMO
	if c.PostForm("username") == "demo" {
		session := sessions.Default(c)
		session.Set("user", "admin")
		session.Set("role", "admin")
		session.Save()
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}
	if stored == "" || captchaInput != stored {
		newCaptcha := utils.GenerateCaptcha()
		parts := strings.Split(newCaptcha, "|")

		session.Set("captcha_answer", parts[0])
		session.Save()

		c.HTML(http.StatusOK, "home.html", gin.H{
			"error":           "Invalid captcha. Please try again.",
			"captchaQuestion": parts[1],
		})
		return
	}

	// --- Load DB ---
	data, err := os.ReadFile(userFile)
	if err != nil {
		c.String(http.StatusInternalServerError, "DB read error")
		return
	}

	var db models.DB
	json.Unmarshal(data, &db)

	// --- Find user by username ONLY ---
	for _, u := range db.Users {

		if u.Username != username {
			continue
		}

		// 🔐 Password check
		if !security.CheckPassword(password, u.Password) {
			security.RecordLoginFailure(username)
			break
		}

		// 🔐 Role validation AFTER password
		if u.Role != role {
			break
		}
		// 🚫 BLOCK INACTIVE EMPLOYEES
		if u.Role == "employee" {
			active := false

			for _, e := range db.Employees {
				if e.EmployeeID == u.Username {
					active = e.IsActive
					break
				}
			}

			if !active {
				newCaptcha := utils.GenerateCaptcha()
				parts := strings.Split(newCaptcha, "|")

				session.Set("captcha_answer", parts[0])
				session.Save()

				c.HTML(http.StatusForbidden, "home.html", gin.H{
					"error":           "Your account is inactive. Please contact HR.",
					"captchaQuestion": parts[1],
				})
				return
			}
		}

		// ✅ SUCCESS
		security.ResetLoginFailures(username)

		session.Set("user", u.Username)
		session.Set("role", u.Role)
		session.Save()

		switch u.Role {
		case "admin":
			c.Redirect(http.StatusFound, "/admin/dashboard")
		case "hr":
			c.Redirect(http.StatusFound, "/admin/employees")
		case "employee":
			c.Redirect(http.StatusFound, "/employee/dashboard")
		case "mentor":
			c.Redirect(http.StatusFound, "/mentor/training")
		case "student":
			c.Redirect(http.StatusFound, "/student/training")
		case "client":
			c.Redirect(http.StatusFound, "/client/dashboard")
		default:
			c.Redirect(http.StatusFound, "/")
		}
		return
	}

	// --- Failure ---
	newCaptcha := utils.GenerateCaptcha()
	parts := strings.Split(newCaptcha, "|")
	session.Set("captcha_answer", parts[0])
	session.Save()

	c.HTML(http.StatusOK, "home.html", gin.H{
		"error":           "Invalid username, password, or role.",
		"captchaQuestion": parts[1],
	})
}

// Logout clears session
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
