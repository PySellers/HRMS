package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"pysellers-erp-go/models"
	"pysellers-erp-go/utils"
)

// ================================
// EMPLOYEE PAYROLL LIST
// ================================
func EmployeePayrollPage(c *gin.Context) {
	session := sessions.Default(c)
	empID := session.Get("user").(string)

	db, _ := utils.ReadDB()

	var myPayrolls []models.Payroll
	for _, p := range db.Payrolls {
		if p.EmployeeID == empID {
			myPayrolls = append(myPayrolls, p)
		}
	}

	c.HTML(http.StatusOK, "employee_payroll.html", gin.H{
		"payrolls": myPayrolls,
	})
}

// ================================
// EMPLOYEE VIEW PAYSLIP
// ================================
func EmployeeViewPayslip(c *gin.Context) {
	session := sessions.Default(c)
	empID := session.Get("user").(string)
	month := c.Param("month")

	db, _ := utils.ReadDB()

	for _, p := range db.Payrolls {
		if p.EmployeeID == empID && p.Month == month {
			c.HTML(http.StatusOK, "employee_payslip_view.html", gin.H{
				"payroll": p,
			})
			return
		}
	}

	c.String(http.StatusNotFound, "Payslip not found")
}

// ================================
// EMPLOYEE DOWNLOAD PAYSLIP
// ================================
func EmployeeDownloadPayslip(c *gin.Context) {
	session := sessions.Default(c)
	empID := session.Get("user").(string)
	month := c.Param("month")

	db, _ := utils.ReadDB()

	var payroll models.Payroll
	var employee models.Employee

	for _, e := range db.Employees {
		if e.EmployeeID == empID {
			employee = e
			break
		}
	}

	for _, p := range db.Payrolls {
		if p.EmployeeID == empID && p.Month == month {
			payroll = p
			break
		}
	}

	path, err := utils.GeneratePayslipPDF(payroll, employee)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	c.FileAttachment(path, "Payslip_"+empID+"_"+month+".pdf")
}
