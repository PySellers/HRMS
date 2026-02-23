package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

// ✅ Use same shared DB file as employee_portal.go
//var dbFile = "data/db.json"

// ShowLeaveRequests — Admin view of all leave requests
func ShowLeaveRequests(c *gin.Context) {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading DB: %v", err)
		return
	}

	var db models.DB
	if err := json.Unmarshal(data, &db); err != nil {
		c.String(http.StatusInternalServerError, "Error parsing DB: %v", err)
		return
	}

	// Combine leave details with employee names for UI
	type LeaveView struct {
		ID         int
		EmployeeID int
		Name       string
		FromDate   string
		ToDate     string
		Reason     string
		Status     string
	}

	var leaveList []LeaveView
	for _, l := range db.Leaves {
		name := "Unknown"
		for _, e := range db.Employees {
			if e.ID == l.EmployeeID {
				name = e.Name
				break
			}
		}
		leaveList = append(leaveList, LeaveView{
			ID:         l.ID,
			EmployeeID: l.EmployeeID,
			Name:       name,
			FromDate:   l.FromDate,
			ToDate:     l.ToDate,
			Reason:     l.Reason,
			Status:     l.Status,
		})
	}

	// Render enhanced UI
	c.HTML(http.StatusOK, "leaves_admin.html", gin.H{
		"leaves": leaveList,
	})
}

// UpdateLeaveStatus — Approve or Reject a leave
func UpdateLeaveStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	action := c.Param("action")

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	json.Unmarshal(data, &db)

	for i, l := range db.Leaves {
		if l.ID == id {
			switch action {
			case "approve":
				db.Leaves[i].Status = "approved"
			case "reject":
				db.Leaves[i].Status = "rejected"
			}
			break
		}
	}

	out, _ := json.MarshalIndent(db, "", "  ")
	_ = os.WriteFile(dbFile, out, 0644)

	c.Redirect(http.StatusFound, "/admin/leaves")
}
