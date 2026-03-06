package main

import (
	"encoding/gob"
	"html/template"
	"net/http"
	"strings"
	"time"
	"os"
	
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
	gin.SetMode(gin.ReleaseMode)

    r := gin.Default()
    r.SetTrustedProxies(nil)

	// ================================
	// STATIC FILES
	// ================================
	r.Static("/static", "./static")

	// ================================
	// TEMPLATE FUNCTIONS
	// ================================
	tmpl := template.Must(
		template.New("").
			Funcs(template.FuncMap{
				"join": strings.Join,
				"title": func(s string) string {
					if len(s) == 0 {
						return s
					}
					return strings.ToUpper(s[:1]) + s[1:]
				},
				"upper": strings.ToUpper,
				"lower": strings.ToLower,
				"add": func(a, b int) int {
					return a + b
				},
			}).
			ParseGlob("templates/*"),
	)

	r.SetHTMLTemplate(tmpl)

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
		admin.GET("/attendance", handlers.AdminAttendance)

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
		admin.GET("/payroll/payslip/:employeeId/:month", handlers.DownloadPayslip)
		admin.GET("/payroll", handlers.PayrollPage)
		admin.POST("/payroll/salary-structure", handlers.SaveSalaryStructure)
		admin.POST("/payroll/generate", handlers.GenerateMonthlyPayroll)
		admin.GET("/employeedetails", handlers.ShowEmployeeDetailsForm)
		admin.POST("/employeedetails/add", handlers.SaveEmployeeDetails)
		admin.GET("/employeedetails/resume/:empID", handlers.DownloadEmployeeResume)
		admin.GET("/payrollmanagement", handlers.PayrollManagement)
		admin.GET("/payrollmanagement/export", handlers.ExportPayrollCSV)
	}

	// ================================
	// EMPLOYEE ROUTES
	// ================================
	employee := r.Group("/employee")
	employee.Use(middleware.RequireLogin(), middleware.RequireRole("employee", "hr", "mentor"))
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
		employee.GET("/payroll", handlers.EmployeePayrollPage)
		employee.GET("/payroll/view/:month", handlers.EmployeeViewPayslip)
		employee.GET("/payroll/download/:month", handlers.EmployeeDownloadPayslip)
	}

	// ================================
	// HR ROUTES (PAYROLL)
	// ================================
	hr := r.Group("/hr")
	hr.Use(middleware.RequireLogin(), middleware.RequireRole("hr"))
	{
		hr.GET("/payroll", handlers.HRPayrollPage)
		hr.GET("/payroll/view/:employeeId/:month", handlers.ViewPayslip)
		hr.GET("/payroll/download/:employeeId/:month", handlers.DownloadPayslip)
		hr.GET("/employeedetails", handlers.ShowEmployeeDetailsForm)
		hr.POST("/employeedetails/add", handlers.SaveEmployeeDetails)
		hr.GET("/employeedetails/resume/:empID", handlers.DownloadEmployeeResume)
		hr.GET("/payrollmanagement", handlers.PayrollManagement)
		hr.GET("/payrollmanagement/export", handlers.ExportPayrollCSV)
	}

	// ================================
	// TRAINING
	// ================================
	training := r.Group("/training")
	training.Use(middleware.RequireLogin())
	{
		// Mentor routes
		mentor := training.Group("/mentor")
		mentor.Use(middleware.RequireRole("mentor"))
		{
			mentor.GET("/", handlers.ShowMentorTraining)
			mentor.POST("/students/add", handlers.AddStudent)
			mentor.POST("/sessions/add", handlers.AddSession)
			mentor.POST("/assignments/create", handlers.CreateAssignment)
			mentor.POST("/attendance/mark", handlers.MarkAttendance)
			mentor.POST("/remarks/add", handlers.AddStudentRemarks)
			mentor.GET("/batches/init", handlers.InitSampleBatches)
			mentor.POST("/batches/add", handlers.AddBatch)
			mentor.GET("/batches/all", handlers.GetAllBatches)
			mentor.GET("/sessions/:id/students", handlers.GetSessionStudents)
			mentor.POST("/sessions/:id/complete", handlers.MarkSessionComplete)
			mentor.POST("/sessions/:id/notes", handlers.AddSessionNotes)
			mentor.GET("/assignments/:id/submissions", handlers.GetAssignmentSubmissions)
			mentor.POST("/assignments/grade", handlers.GradeAssignment)
			mentor.GET("/students/:id", handlers.GetStudentDetails)
		}

		// Student routes
		student := training.Group("/student")
		student.Use(middleware.RequireRole("student"))
		{
			student.GET("/", handlers.ShowStudentTraining)
			student.POST("/assignment/submit", handlers.UploadSubmission)
			student.GET("/assignment/:id", handlers.GetAssignmentDetails)
			student.GET("/session/:id", handlers.GetSessionDetails)
			student.POST("/certificate/request", handlers.RequestCertificate)
		}
	}

	// Redirect routes for backward compatibility
	r.GET("/mentor/training", handlers.RedirectMentorTraining)
	r.GET("/student/training", handlers.RedirectStudentTraining)

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
	port := "8080"
if p := os.Getenv("PORT"); p != "" {
    port = p
}

r.Run(":" + port)
}