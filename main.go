package main

import (
	"encoding/gob"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"pysellers-erp-go/handlers"
	"pysellers-erp-go/middleware"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(map[string]interface{}{})
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// ================================
	// STATIC FILES
	// ================================
	r.Static("/static", "./static")

	// ================================
	// TEMPLATE FUNCTIONS
	// ================================
	r.SetFuncMap(template.FuncMap{
		"join": strings.Join,
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	})

	r.LoadHTMLGlob("templates/*")

	// ================================
	// SESSION SETUP
	// ================================
	store := cookie.NewStore([]byte("superstrong_random_secret_key_32chars"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   false, // true in HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("mysession", store))

	// ================================
	// DEMO LOGIN
	// ================================
	r.GET("/demo/admin", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user", "admin")
		session.Set("role", "admin")
		session.Save()
		c.Redirect(http.StatusFound, "/admin/dashboard")
	})

	r.GET("/demo/employee", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user", "employee1")
		session.Set("role", "employee")
		session.Save()
		c.Redirect(http.StatusFound, "/employee/dashboard")
	})

	// ================================
	// KEEP ALIVE
	// ================================
	r.GET("/ping", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("lastActivity", time.Now())
		session.Save()
		c.JSON(200, gin.H{"status": "alive"})
	})

	// ================================
	// PUBLIC ROUTES
	// ================================
	r.GET("/", handlers.ShowHome)
	r.POST("/login", handlers.Login)
	r.GET("/logout", handlers.Logout)

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", nil)
	})

	r.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.html", nil)
	})

	// ================================
	// ADMIN / HR ROUTES
	// ================================
	admin := r.Group("/admin")
	admin.Use(middleware.RequireLogin(), middleware.RequireRole("admin", "hr"))
	{
		admin.GET("/dashboard", handlers.ShowAdminDashboard)

		// EMPLOYEES
		admin.GET("/employees", handlers.ShowEmployees)
		admin.POST("/employees/add", handlers.AddEmployee)
		admin.POST("/employees/toggle/:id", handlers.ToggleEmployeeStatus)
		admin.POST("/employees/update/:id", handlers.UpdateEmployee)
		admin.GET("/employees/delete/:id", handlers.DeleteEmployee)

		// LEAVES
		admin.GET("/leaves", handlers.ShowLeaveRequests)
		admin.GET("/leaves/:id/:action", handlers.UpdateLeaveStatus)

		// FINANCE
		admin.GET("/finance", handlers.ShowFinanceDashboard)
		admin.POST("/finance/add", handlers.AddFinance)

		// NEWS
		admin.GET("/news", handlers.ShowNewsPage)
		admin.POST("/news/add", handlers.AddNews)

		// HR REQUESTS
		admin.GET("/hrrequests", handlers.ShowHRRequests)
		admin.GET("/hrrequests/:id/:action", handlers.UpdateHRRequestStatus)

		// POLICIES
		admin.GET("/policies", handlers.ShowPolicies)
		admin.POST("/policies/add", handlers.AddPolicy)

		// MARKETING
		admin.GET("/marketing", handlers.ShowMarketingPage)
		admin.POST("/marketing/add", handlers.AddMarketingActivity)
		admin.GET("/marketing/weekly-summary", handlers.GenerateWeeklySummary)

		// RECRUITMENT
		admin.GET("/recruitment", handlers.ShowRecruitments)
		admin.POST("/recruitment/add", handlers.AddRecruitment)
		admin.GET("/recruitment/:id/:action", handlers.UpdateRecruitmentStatus)

		// TOOLS
		admin.GET("/tools", handlers.ShowAdminTools)
		admin.POST("/backup", handlers.BackupData)
		admin.GET("/export", handlers.ExportData)
		admin.POST("/clear", handlers.ClearDummyData)
		admin.GET("/health", handlers.SystemHealth)

		// SETTINGS
		admin.GET("/settings", handlers.ShowSettings)
		admin.POST("/settings/update", handlers.UpdateSettings)
	}

	// ================================
	// EMPLOYEE ROUTES
	// ================================
	employee := r.Group("/employee")
	employee.Use(middleware.RequireLogin(), middleware.RequireRole("employee"))
	{
		employee.GET("/dashboard", handlers.ShowEmployeeDashboard)
		employee.POST("/profile", handlers.UpdateProfile)
		employee.POST("/timein", handlers.TimeIn)
		employee.POST("/timeout", handlers.TimeOut)
		employee.POST("/leave", handlers.ApplyLeave)
		employee.POST("/skills", handlers.AddSkillCertification)
		employee.POST("/uploadpic", handlers.UploadProfilePic)
		employee.GET("/change-password", handlers.ShowChangePassword)
		employee.POST("/change-password", handlers.ChangePassword)

		employee.GET("/hrrequest", handlers.ShowEmployeeHRForm)
		employee.POST("/hrrequest", handlers.SubmitHRRequest)
	}

	// ================================
	// TRAINING
	// ================================
	training := r.Group("/training")
	training.Use(middleware.RequireLogin(), middleware.RequireRole("mentor", "student"))
	{
		training.GET("/mentor", handlers.ShowMentorTraining)
		training.GET("/student", handlers.ShowStudentTraining)
	}

	// ================================
	// CLIENT
	// ================================
	client := r.Group("/client")
	client.Use(middleware.RequireLogin(), middleware.RequireRole("client"))
	{
		client.GET("/dashboard", handlers.ShowClientDashboard)
	}

	// ================================
	// START SERVER
	// ================================
	r.Run(":9090")
}
