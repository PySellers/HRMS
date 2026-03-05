package handlers

import (
	"encoding/json"
	"log"
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

// =======================
// SHOW LOGIN PAGE
// =======================
func ShowHome(c *gin.Context) {
	session := sessions.Default(c)

	captcha := utils.GenerateCaptcha()
	parts := strings.Split(captcha, "|")
	if len(parts) != 2 {
		c.String(http.StatusInternalServerError, "Captcha error")
		return
	}

	session.Set("captcha_answer", parts[0])
	session.Save()

	c.HTML(http.StatusOK, "home.html", gin.H{
		"captchaQuestion": parts[1],
	})
}

// =======================
// LOGIN HANDLER
// =======================
// Login - Handles user login
func Login(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")
	activeRole := strings.TrimSpace(c.PostForm("role"))
	captchaInput := strings.TrimSpace(c.PostForm("captcha"))

	// Add debug logs
	log.Printf("🔐 Login attempt - Username/Email: %s, Role: %s", username, activeRole)

	session := sessions.Default(c)

	// -------- CAPTCHA --------
	stored, _ := session.Get("captcha_answer").(string)
	if stored == "" || captchaInput != stored {
		log.Printf("❌ Captcha failed - Expected: %s, Got: %s", stored, captchaInput)
		captcha := utils.GenerateCaptcha()
		parts := strings.Split(captcha, "|")

		session.Set("captcha_answer", parts[0])
		session.Save()

		c.HTML(http.StatusOK, "home.html", gin.H{
			"error":           "Invalid captcha",
			"captchaQuestion": parts[1],
		})
		return
	}

	// -------- LOAD DB --------
	data, err := os.ReadFile(userFile)
	if err != nil {
		log.Printf("❌ DB read error: %v", err)
		c.String(http.StatusInternalServerError, "DB read error")
		return
	}

	var db models.DB
	json.Unmarshal(data, &db)

	// Log all users in DB for debugging
	log.Printf("📋 Total users in DB: %d", len(db.Users))
	for i, u := range db.Users {
		log.Printf("  User %d: Username=%s, Email=%s, Role=%s", i+1, u.Username, u.Email, u.Role)
	}

	// -------- AUTH --------
	for _, u := range db.Users {
		log.Printf("🔍 Checking user: Username=%s, Email=%s against input=%s", u.Username, u.Email, username)

		if u.Username != username && u.Email != username {
			continue
		}
		log.Printf("✅ User found: %s", u.Username)

		// Password check
		log.Printf("🔑 Checking password for user %s", u.Username)
		if !security.CheckPassword(password, u.Password) {
			log.Printf("❌ Password mismatch for user %s", u.Username)
			security.RecordLoginFailure(username)
			break
		}
		log.Printf("✅ Password correct for user %s", u.Username)

		// -------- ROLE DELEGATION RULES --------
		allowed := false
		switch u.Role {
		case "admin":
			allowed = true
		case "hr":
			allowed = activeRole == "hr" || activeRole == "employee"
		case "mentor":
			allowed = activeRole == "mentor" || activeRole == "employee"
		case "employee":
			allowed = activeRole == "employee"
		case "student":
			allowed = activeRole == "student"
		case "client":
			allowed = activeRole == "client"
		}

		log.Printf("🎭 Role check - User role: %s, Selected role: %s, Allowed: %v", u.Role, activeRole, allowed)

		if !allowed {
			log.Printf("❌ Role not allowed for user %s", u.Username)
			break
		}

		// -------- BLOCK INACTIVE EMPLOYEES --------
		if activeRole == "employee" {
			active := false
			for _, e := range db.Employees {
				if e.EmployeeID == u.Username {
					active = e.IsActive
					break
				}
			}
			if !active {
				log.Printf("❌ Employee account inactive: %s", u.Username)
				c.HTML(http.StatusForbidden, "home.html", gin.H{
					"error": "Your account is inactive. Contact HR.",
				})
				return
			}
		}

		// -------- SUCCESS --------
		security.ResetLoginFailures(username)

		session.Set("user", u.Username)
		session.Set("db_role", u.Role)  // REAL ROLE
		session.Set("role", activeRole) // ACTIVE ROLE
		session.Save()

		log.Printf("✅ Login successful for %s as %s", u.Username, activeRole)

		// -------- REDIRECT --------
		switch activeRole {
		case "admin":
			c.Redirect(http.StatusFound, "/admin/dashboard")
		case "hr":
			c.Redirect(http.StatusFound, "/admin/employees")
		case "employee":
			c.Redirect(http.StatusFound, "/employee/dashboard")
		case "mentor":
			c.Redirect(http.StatusFound, "/training/mentor")
		case "student":
			c.Redirect(http.StatusFound, "/training/student")
		case "client":
			c.Redirect(http.StatusFound, "/client/dashboard")
		default:
			c.Redirect(http.StatusFound, "/")
		}
		return
	}

	// -------- FAILURE --------
	log.Printf("❌ Login failed for user: %s", username)
	captcha := utils.GenerateCaptcha()
	parts := strings.Split(captcha, "|")
	session.Set("captcha_answer", parts[0])
	session.Save()

	c.HTML(http.StatusOK, "home.html", gin.H{
		"error":           "Invalid username or password",
		"captchaQuestion": parts[1],
	})
}

// =======================
// LOGOUT
// =======================
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
