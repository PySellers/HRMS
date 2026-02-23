package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

var reportDB = "data/db.json"

func ShowReports(c *gin.Context) {
	data, _ := os.ReadFile(reportDB)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	// ---- Compute attendance count per employee ----
	type AttendanceSummary struct {
		Name        string
		DaysPresent int
	}
	var attendance []AttendanceSummary
	for _, e := range db.Employees {
		count := 0
		for _, a := range db.Attendance {
			if a.EmployeeID == e.ID && a.TimeIn != "" {
				count++
			}
		}
		attendance = append(attendance, AttendanceSummary{Name: e.Name, DaysPresent: count})
	}

	// ---- Leave summary ----
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

	c.HTML(http.StatusOK, "reports.html", gin.H{
		"attendance": attendance,
		"stats":      stats,
	})
}
