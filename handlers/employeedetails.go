package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pysellers-erp-go/models"
	"pysellers-erp-go/utils"

	"github.com/gin-gonic/gin"
)


// =====================================================
// SHOW EMPLOYEE DETAILS FORM (HR / ADMIN)
// =====================================================
func ShowEmployeeDetailsForm(c *gin.Context) {

	empID := c.Query("emp_id")
	db, _ := utils.ReadDB()

	// ✅ ALWAYS pass concrete struct
	employee := models.EmployeeDetails{}

	// Load existing HR details if available
	if empID != "" {
		for _, d := range db.EmployeeDetails {
			if d.EmployeeID == empID {
				employee = d
				break
			}
		}
	}

	c.HTML(http.StatusOK, "employeedetails.html", gin.H{
		"employee": employee,
	})
}


// =====================================================
// SAVE / UPDATE EMPLOYEE FULL DETAILS (WITH RESUME)
// =====================================================
func SaveEmployeeDetails(c *gin.Context) {

	utils.DBMutex.Lock()
	defer utils.DBMutex.Unlock()

	db, _ := utils.ReadDB()
	employeeID := c.PostForm("employee_id")

	// =========================
	// HANDLE RESUME UPLOAD
	// =========================
	var resumePath string

	file, err := c.FormFile("resume")
	if err == nil && file != nil {

		// Size validation (2MB)
		if file.Size > 2*1024*1024 {
			c.String(http.StatusBadRequest, "Resume must be under 2MB")
			return
		}

		// Extension validation
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
			c.String(http.StatusBadRequest, "Only PDF, DOC, DOCX allowed")
			return
		}

		uploadDir := "static/uploads/resumes"
		_ = os.MkdirAll(uploadDir, os.ModePerm)

		// Auto-replace existing resume
		resumePath = filepath.Join(uploadDir, employeeID+ext)
		_ = os.Remove(resumePath)

		if err := c.SaveUploadedFile(file, resumePath); err != nil {
			c.String(http.StatusInternalServerError, "Resume upload failed")
			return
		}
	}

	// =========================
	// BUILD EMPLOYEE DETAILS
	// =========================
	details := models.EmployeeDetails{
		EmployeeID:  employeeID,
		FullName:    c.PostForm("full_name"),
		Phone:       c.PostForm("phone"),
		Email:       c.PostForm("email"),
		Address:     c.PostForm("address"),
		DateOfBirth: c.PostForm("dob"),
		Gender:      c.PostForm("gender"),

		Designation:    c.PostForm("designation"),
		Department:     c.PostForm("department"),
		DateOfJoining:  c.PostForm("doj"),
		EmploymentType: c.PostForm("employment_type"),

		PAN:     c.PostForm("pan"),
		Aadhaar: c.PostForm("aadhaar"),

		BankName:      c.PostForm("bank_name"),
		AccountNumber: c.PostForm("account_number"),
		IFSC:          c.PostForm("ifsc"),

		Degree:     c.PostForm("degree"),
		Experience: c.PostForm("experience"),

		Resume:    resumePath,
		UpdatedAt: time.Now().Format("2006-01-02 15:04"),
	}

	// =========================
	// UPDATE IF EXISTS
	// =========================
	for i, d := range db.EmployeeDetails {
		if d.EmployeeID == employeeID {

			// Preserve old resume if no new upload
			if resumePath == "" {
				details.Resume = d.Resume
			}

			details.CreatedAt = d.CreatedAt
			db.EmployeeDetails[i] = details

			utils.WriteDB(db)
			c.Redirect(http.StatusFound, "/admin/employees")
			return
		}
	}

	// =========================
	// CREATE NEW
	// =========================
	details.CreatedAt = time.Now().Format("2006-01-02 15:04")
	db.EmployeeDetails = append(db.EmployeeDetails, details)

	utils.WriteDB(db)
	c.Redirect(http.StatusFound, "/admin/employees")
}


// =====================================================
// DOWNLOAD EMPLOYEE RESUME (FIXES 404)
// =====================================================
func DownloadEmployeeResume(c *gin.Context) {

	empID := c.Param("empID")
	db, _ := utils.ReadDB()

	for _, d := range db.EmployeeDetails {
		if d.EmployeeID == empID && d.Resume != "" {

			// Extra safety: file must exist
			if _, err := os.Stat(d.Resume); err != nil {
				c.String(http.StatusNotFound, "Resume file missing")
				return
			}

			c.FileAttachment(d.Resume, filepath.Base(d.Resume))
			return
		}
	}

	c.String(http.StatusNotFound, "Resume not found")
}