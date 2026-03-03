package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"pysellers-erp-go/models"

	"github.com/gin-gonic/gin"
)

/*
================================

	ANALYTICS VIEW MODEL

================================
*/
type AttendanceAnalytics struct {
	Date      string  `json:"date"`
	Employee  string  `json:"employee"`
	Hours     float64 `json:"hours"`
	TargetMet bool    `json:"targetMet"`
}

func AdminAttendance(c *gin.Context) {

	// ===============================
	// LOAD DATABASE
	// ===============================
	data, err := os.ReadFile(dbFile)
	if err != nil {
		c.String(500, "Failed to read database")
		return
	}

	var db models.DB
	_ = json.Unmarshal(data, &db)

	// ===============================
	// EMPLOYEE ID → NAME MAP
	// ===============================
	empMap := make(map[int]string)
	for _, e := range db.Employees {
		empMap[e.ID] = e.Name
	}

	// ===============================
	// SORT ATTENDANCE BY DATE DESC
	// ===============================
	sort.Slice(db.Attendance, func(i, j int) bool {
		return db.Attendance[i].Date > db.Attendance[j].Date
	})

	// ===============================
	// VIEW MODELS
	// ===============================
	type SessionView struct {
		TimeIn       string
		TimeOut      string
		EmployeeName string
	}

	type AttendanceView struct {
		Date      string
		TotalTime string
		Sessions  []SessionView
	}

	var viewData []AttendanceView
	var analytics []AttendanceAnalytics

	// ===============================
	// BUILD VIEW + ANALYTICS
	// ===============================
	for _, a := range db.Attendance {

		// Sort sessions (latest first)
		sort.Slice(a.Sessions, func(i, j int) bool {
			return a.Sessions[i].TimeIn > a.Sessions[j].TimeIn
		})

		var sessions []SessionView
		for _, s := range a.Sessions {
			sessions = append(sessions, SessionView{
				TimeIn:       s.TimeIn,
				TimeOut:      s.TimeOut,
				EmployeeName: empMap[a.EmployeeID],
			})
		}

		viewData = append(viewData, AttendanceView{
			Date:      a.Date,
			TotalTime: a.TotalTime,
			Sessions:  sessions,
		})

		// ===============================
		// ANALYTICS (8 HOURS CHECK)
		// ===============================
		hours := 0.0

		if a.TotalTime != "" {
			parts := strings.Split(a.TotalTime, ":")
			if len(parts) == 3 {
				h, _ := strconv.Atoi(parts[0])
				m, _ := strconv.Atoi(parts[1])
				s, _ := strconv.Atoi(parts[2])
				hours = float64(h) + float64(m)/60 + float64(s)/3600
			}
		}

		analytics = append(analytics, AttendanceAnalytics{
			Date:      a.Date,
			Employee:  empMap[a.EmployeeID],
			Hours:     round(hours, 2),
			TargetMet: hours >= 8,
		})
	}

	// ===============================
	// CONVERT ANALYTICS → JSON (JS SAFE)
	// ===============================
	analyticsJSON, _ := json.Marshal(analytics)

	// ===============================
	// RENDER TEMPLATE
	// ===============================
	c.HTML(http.StatusOK, "attendance.html", gin.H{
		"attendance": viewData,
		"analytics":  template.JS(analyticsJSON), // 🔥 IMPORTANT
	})
}

// ===============================
// UTILITY: ROUND FLOAT
// ===============================
func round(val float64, precision int) float64 {
	pow := 1.0
	for i := 0; i < precision; i++ {
		pow *= 10
	}
	return float64(int(val*pow+0.5)) / pow
}
