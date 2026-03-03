package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"pysellers-erp-go/models"

	"github.com/gin-gonic/gin"
)

var reportDB = "data/db.json"

func ShowReports(c *gin.Context) {
	data, err := os.ReadFile(reportDB)
	if err != nil {
		c.String(500, "Failed to read database")
		return
	}

	var db models.DB
	_ = json.Unmarshal(data, &db)

	// ================================
	// ATTENDANCE SUMMARY (ENHANCED)
	// ================================
	type AttendanceSummary struct {
		Name        string
		DaysPresent int
	}

	var attendance []AttendanceSummary

	for _, e := range db.Employees {
		daysPresent := 0

		for _, a := range db.Attendance {
			if a.EmployeeID == e.ID && len(a.Sessions) > 0 {
				daysPresent++ // one record = one day
			}
		}

		attendance = append(attendance, AttendanceSummary{
			Name:        e.Name,
			DaysPresent: daysPresent,
		})
	}

	// ================================
	// LEAVE SUMMARY (UNCHANGED)
	// ================================
	totalLeaves := len(db.Leaves)
	pending, approved, rejected := 0, 0, 0

	for _, l := range db.Leaves {
		switch l.Status {
		case "pending":
			pending++
		case "approved":
			approved++
		case "rejected":
			rejected++
		}
	}

	stats := map[string]int{
		"totalEmployees": len(db.Employees),
		"totalLeaves":    totalLeaves,
		"approved":       approved,
		"pending":        pending,
		"rejected":       rejected,
	}

	// ================================
	// RENDER REPORT
	// ================================
	c.HTML(http.StatusOK, "reports.html", gin.H{
		"attendance": attendance,
		"stats":      stats,
	})
}
