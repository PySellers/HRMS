package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pysellers-erp-go/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const dbFile = "data/db.json"

// ================================
// EMPLOYEE DASHBOARD
// ================================
func EmployeeDashboard(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user") // EmployeeID (e.g. PY-8017)

	if username == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	data, err := os.ReadFile(dbFile)
	if err != nil {
		c.String(500, "Failed to read database")
		return
	}

	var db models.DB
	_ = json.Unmarshal(data, &db)

	var emp models.Employee
	var attendance []models.Attendance
	var leaves []models.Leave
	var payroll []models.Payroll

	// ✅ MATCH USING EmployeeID
	for _, e := range db.Employees {
		if strings.EqualFold(e.EmployeeID, username.(string)) {
			emp = e
			break
		}
	}

	if emp.ID == 0 {
		c.HTML(http.StatusOK, "employee_dashboard.html", gin.H{
			"error": "Employee profile not found. Contact admin.",
		})
		return
	}

	for _, a := range db.Attendance {
		if a.EmployeeID == emp.ID {
			attendance = append(attendance, a)
		}
	}

	for _, l := range db.Leaves {
		if l.EmployeeID == emp.ID {
			leaves = append(leaves, l)
		}
	}

	for _, p := range db.Payroll {
		if p.EmployeeID == emp.ID {
			payroll = append(payroll, p)
		}
	}

	c.HTML(http.StatusOK, "employee_dashboard.html", gin.H{
		"employee":   emp,
		"attendance": attendance,
		"leaves":     leaves,
		"payroll":    payroll,
		"news":       db.News,
		"policies":   db.Policies,
		"hrrequests": filterHRRequests(emp.ID, db.HRRequests),
	})
}

// ================================
// UPDATE PROFILE
// ================================
func UpdateProfile(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	for i, e := range db.Employees {
		if strings.EqualFold(e.EmployeeID, username.(string)) {
			db.Employees[i].Phone = c.PostForm("phone")
			db.Employees[i].Address = c.PostForm("address")
			break
		}
	}

	saveDB(db)
	c.Redirect(http.StatusFound, "/employee/dashboard")
}

// ================================
// TIME IN
// ================================
func TimeIn(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	if username == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	empID := getEmployeeID(username.(string), db)
	if empID == -1 {
		c.String(401, "Employee not found")
		return
	}

	today := time.Now().Format("2006-01-02")

	for _, a := range db.Attendance {
		if a.EmployeeID == empID && a.Date == today && a.TimeIn != "" {
			c.Redirect(http.StatusFound, "/employee/dashboard")
			return
		}
	}

	db.Attendance = append(db.Attendance, models.Attendance{
		ID:         len(db.Attendance) + 1,
		EmployeeID: empID,
		Date:       today,
		TimeIn:     time.Now().Format("15:04:05"),
		Latitude:   c.PostForm("latitude"),
		Longitude:  c.PostForm("longitude"),
	})

	saveDB(db)
	c.Redirect(http.StatusFound, "/employee/dashboard")
}

// ================================
// TIME OUT
// ================================
func TimeOut(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	empID := getEmployeeID(username.(string), db)

	for i, a := range db.Attendance {
		if a.EmployeeID == empID && a.Date == time.Now().Format("2006-01-02") {
			db.Attendance[i].TimeOut = time.Now().Format("15:04:05")
			break
		}
	}

	saveDB(db)
	c.Redirect(http.StatusFound, "/employee/dashboard")
}

// ================================
// APPLY LEAVE
// ================================
func ApplyLeave(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	empID := getEmployeeID(username.(string), db)
	if empID == -1 {
		c.String(401, "Employee not found")
		return
	}

	db.Leaves = append(db.Leaves, models.Leave{
		ID:         len(db.Leaves) + 1,
		EmployeeID: empID,
		FromDate:   c.PostForm("from_date"),
		ToDate:     c.PostForm("to_date"),
		Reason:     c.PostForm("reason"),
		Status:     "pending",
	})

	saveDB(db)
	c.Redirect(http.StatusFound, "/employee/dashboard")
}

// ================================
// ADD SKILL / CERTIFICATION
// ================================
func AddSkillCertification(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	skill := c.PostForm("skill")
	cert := c.PostForm("certification")

	for i, e := range db.Employees {
		if strings.EqualFold(e.EmployeeID, username.(string)) {
			if skill != "" {
				db.Employees[i].Skills = append(db.Employees[i].Skills, skill)
			}
			if cert != "" {
				db.Employees[i].Certifications = append(db.Employees[i].Certifications, cert)
			}
			break
		}
	}

	saveDB(db)
	c.Redirect(http.StatusFound, "/employee/dashboard")
}

// ================================
// UPLOAD PROFILE PIC
// ================================
func UploadProfilePic(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	file, err := c.FormFile("profilepic")
	if err != nil {
		c.String(400, "Upload failed")
		return
	}

	filename := filepath.Base(file.Filename)
	savePath := filepath.Join("static/profile_pics", filename)
	_ = c.SaveUploadedFile(file, savePath)

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	for i, e := range db.Employees {
		if strings.EqualFold(e.EmployeeID, username.(string)) {
			db.Employees[i].ProfilePic = filename
			break
		}
	}

	saveDB(db)
	c.Redirect(http.StatusFound, "/employee/dashboard")
}

// ================================
// CHANGE PASSWORD
// ================================
func ShowChangePassword(c *gin.Context) {
	c.HTML(http.StatusOK, "change_password.html", nil)
}

func ChangePassword(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	oldPwd := c.PostForm("old_password")
	newPwd := c.PostForm("new_password")
	confirm := c.PostForm("confirm_password")

	if newPwd != confirm {
		c.String(400, "Passwords do not match")
		return
	}

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	for i, u := range db.Users {
		if u.Username == username && u.Password == oldPwd {
			db.Users[i].Password = newPwd
			saveDB(db)
			c.Redirect(http.StatusFound, "/employee/dashboard")
			return
		}
	}

	c.String(401, "Old password incorrect")
}

// ================================
// HELPERS
// ================================
func getEmployeeID(username string, db models.DB) int {
	for _, e := range db.Employees {
		if strings.EqualFold(e.EmployeeID, username) {
			return e.ID
		}
	}
	return -1
}

func filterHRRequests(empID int, all []models.HRRequest) []models.HRRequest {
	var result []models.HRRequest
	for _, r := range all {
		if r.EmployeeID == empID {
			result = append(result, r)
		}
	}
	return result
}

func saveDB(db models.DB) {
	data, _ := json.MarshalIndent(db, "", "  ")
	_ = os.WriteFile(dbFile, data, 0644)
}

// ShowEmployeeDashboard is a route alias for EmployeeDashboard
// (kept to avoid breaking main.go)
func ShowEmployeeDashboard(c *gin.Context) {
	EmployeeDashboard(c)
}
