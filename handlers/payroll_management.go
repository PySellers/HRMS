package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pysellers-erp-go/models"

	"pysellers-erp-go/utils"

	"github.com/gin-gonic/gin"
)

type DayAttendance struct {
	Date     string
	Status   string // PRESENT, ABSENT, PAID_LEAVE, UNPAID_LEAVE
	WorkTime string
	Sessions []models.AttendanceSession
}

type EmployeePayrollView struct {
	EmployeeID   string
	Name         string
	Days         []DayAttendance
	WorkingDays  int
	PresentDays  int
	PaidLeaves   int
	UnpaidLeaves int
	Deduction    float64
	NetSalary    float64
}

func PayrollManagement(c *gin.Context) {

	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	db, _ := utils.ReadDB()

	start, _ := time.Parse("2006-01", month)
	daysInMonth := start.AddDate(0, 1, -1).Day()

	var view []EmployeePayrollView

	for _, emp := range db.Employees {

		if emp.Role != "employee" {
			continue
		}

		row := EmployeePayrollView{
			EmployeeID: emp.EmployeeID,
			Name:       emp.Name,
		}

		for d := 1; d <= daysInMonth; d++ {
			date := time.Date(start.Year(), start.Month(), d, 0, 0, 0, 0, time.Local)
			dateStr := date.Format("2006-01-02")

			day := DayAttendance{Date: dateStr, Status: "ABSENT"}

			// Attendance check
			for _, a := range db.Attendance {
				if a.EmployeeID == emp.ID && a.Date == dateStr {
					day.Status = "PRESENT"
					day.WorkTime = a.TotalTime
					day.Sessions = a.Sessions
					row.PresentDays++
				}
			}

			// Leave check
			for _, l := range db.Leaves {
				if l.EmployeeID == emp.ID && l.Status == "approved" {
					if dateStr >= l.FromDate && dateStr <= l.ToDate {
						if row.PaidLeaves < 2 {
							day.Status = "PAID_LEAVE"
							row.PaidLeaves++
						} else {
							day.Status = "UNPAID_LEAVE"
							row.UnpaidLeaves++
						}
					}
				}
			}

			row.Days = append(row.Days, day)
			if day.Status != "ABSENT" {
				row.WorkingDays++
			}
		}

		// Salary
		for _, p := range db.Payrolls {
			if p.EmployeeID == emp.EmployeeID && strings.HasPrefix(p.Month, month) {
				perDay := p.Basic / float64(daysInMonth)
				row.Deduction = perDay * float64(row.UnpaidLeaves)
				row.NetSalary = p.Net - row.Deduction
			}
		}

		view = append(view, row)
	}

	c.HTML(http.StatusOK, "admin_payroll_management.html", gin.H{
		"Month": month,
		"Data":  view,
	})
}
func ExportPayrollCSV(c *gin.Context) {

	db, _ := utils.ReadDB()

	c.Header("Content-Disposition", "attachment; filename=payroll.csv")
	c.Header("Content-Type", "text/csv")

	w := c.Writer
	w.Write([]byte("EmployeeID,Name,PresentDays,PaidLeaves,UnpaidLeaves,NetSalary\n"))

	for _, emp := range db.Employees {

		if emp.Role != "employee" {
			continue
		}

		var present, paid, unpaid int
		var net float64

		for _, p := range db.Payrolls {
			if p.EmployeeID == emp.EmployeeID {
				net = p.Net
			}
		}

		line := emp.EmployeeID + "," +
			emp.Name + "," +
			strconv.Itoa(present) + "," +
			strconv.Itoa(paid) + "," +
			strconv.Itoa(unpaid) + "," +
			fmt.Sprintf("%.2f", net) + "\n"

		w.Write([]byte(line))
	}

}
