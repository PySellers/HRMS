package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pysellers-erp-go/models"

	"fmt"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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

	// Match using EmployeeID
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
	// ================================
	// MONTHLY SUMMARY CALCULATION
	// ================================
	days, totalDuration := calculateMonthlySummary(attendance)

	monthTotal := fmt.Sprintf("%02d:%02d:%02d",
		int(totalDuration.Hours()),
		int(totalDuration.Minutes())%60,
		int(totalDuration.Seconds())%60,
	)

	for _, l := range db.Leaves {
		if l.EmployeeID == emp.ID {
			leaves = append(leaves, l)
		}
	}

	for _, p := range db.Payrolls {
		if p.EmployeeID == emp.EmployeeID {
			payroll = append(payroll, p)
		}
	}

	c.HTML(http.StatusOK, "employee_dashboard.html", gin.H{
		"employee":    emp,
		"attendance":  attendance,
		"leaves":      leaves,
		"payroll":     payroll,
		"news":        db.News,
		"policies":    db.Policies,
		"hrrequests":  filterHRRequests(emp.ID, db.HRRequests),
		"month_days":  days,
		"month_total": monthTotal,
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
	now := time.Now().Format("15:04:05")

	for i, a := range db.Attendance {
		if a.EmployeeID == empID && a.Date == today {

			// ❌ Prevent Time In if last session not closed
			if len(a.Sessions) > 0 && a.Sessions[len(a.Sessions)-1].TimeOut == "" {
				c.Redirect(http.StatusFound, "/employee/dashboard")
				return
			}

			db.Attendance[i].Sessions = append(db.Attendance[i].Sessions, models.AttendanceSession{
				TimeIn: now,
			})

			saveDB(db)
			c.Redirect(http.StatusFound, "/employee/dashboard")
			return
		}
	}

	// First entry of the day
	db.Attendance = append(db.Attendance, models.Attendance{
		ID:         len(db.Attendance) + 1,
		EmployeeID: empID,
		Date:       today,
		Sessions: []models.AttendanceSession{
			{TimeIn: now},
		},
		TotalTime: "00:00:00",
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
	today := time.Now().Format("2006-01-02")
	now := time.Now()

	for i, a := range db.Attendance {
		if a.EmployeeID == empID && a.Date == today {
			last := len(a.Sessions) - 1
			if last >= 0 && a.Sessions[last].TimeOut == "" {

				db.Attendance[i].Sessions[last].TimeOut = now.Format("15:04:05")
				db.Attendance[i].TotalTime = calculateTotalTime(db.Attendance[i].Sessions)

				saveDB(db)
				break
			}
		}
	}

	c.Redirect(http.StatusFound, "/employee/dashboard")
}
func calculateTotalTime(sessions []models.AttendanceSession) string {
	var total time.Duration

	for _, s := range sessions {
		if s.TimeIn != "" && s.TimeOut != "" {
			in, _ := time.Parse("15:04:05", s.TimeIn)
			out, _ := time.Parse("15:04:05", s.TimeOut)
			total += out.Sub(in)
		}
	}

	h := int(total.Hours())
	m := int(total.Minutes()) % 60
	s := int(total.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
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
func calculateMonthlySummary(att []models.Attendance) (int, time.Duration) {
	var days int
	var total time.Duration
	now := time.Now()

	for _, a := range att {
		d, err := time.Parse("2006-01-02", a.Date)
		if err != nil {
			continue
		}

		if d.Year() == now.Year() && d.Month() == now.Month() {
			tt := a.TotalTime
			if tt == "" {
				tt = "00:00:00"
			}

			parts := strings.Split(tt, ":")
			h, _ := strconv.Atoi(parts[0])
			m, _ := strconv.Atoi(parts[1])
			s, _ := strconv.Atoi(parts[2])

			total += time.Duration(h)*time.Hour +
				time.Duration(m)*time.Minute +
				time.Duration(s)*time.Second
			days++
		}
	}
	return days, total
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
// CHANGE PASSWORD (FIXED)
// ================================
func ShowChangePassword(c *gin.Context) {
	c.HTML(http.StatusOK, "change_password.html", nil)
}

func ChangePassword(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")

	if username == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	userID := username.(string)

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
		if u.Username == userID {

			// ✅ Only check old password if NOT first login
			if !u.ForceChangePassword {
				err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPwd))
				if err != nil {
					c.String(401, "Old password incorrect")
					return
				}
			}

			// Hash new password
			hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)

			db.Users[i].Password = string(hashedPwd)
			db.Users[i].ForceChangePassword = false // ✅ first login completed

			saveDB(db)
			c.Redirect(http.StatusFound, "/employee/dashboard")
			return
		}
	}

	c.String(401, "User not found")
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

// Route alias (do not remove)
func ShowEmployeeDashboard(c *gin.Context) {
	EmployeeDashboard(c)
}
