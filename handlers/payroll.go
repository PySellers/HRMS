package handlers

import (
	"net/http"
	"strconv"
	"time"

	"pysellers-erp-go/models"
	"pysellers-erp-go/utils"

	"github.com/gin-gonic/gin"
)

// ================================
// PAYROLL PAGE (ADMIN + HR)
// ================================
func PayrollPage(c *gin.Context) {

	db, _ := utils.ReadDB()

	// Attach employee name to payroll
	for i := range db.Payrolls {

		for _, emp := range db.Employees {

			if emp.EmployeeID == db.Payrolls[i].EmployeeID {

				db.Payrolls[i].EmployeeName = emp.Name
				break
			}
		}
	}

	c.HTML(http.StatusOK, "payroll.html", gin.H{
		"employees": db.Employees,
		"salaries":  db.SalaryStructures,
		"payrolls":  db.Payrolls,
	})
}

// ================================
// SAVE SALARY STRUCTURE
// ================================
func SaveSalaryStructure(c *gin.Context) {

	utils.DBMutex.Lock()
	defer utils.DBMutex.Unlock()

	employeeID := c.PostForm("employee_id")

	basic, _ := strconv.ParseFloat(c.PostForm("basic"), 64)
	hra, _ := strconv.ParseFloat(c.PostForm("hra"), 64)
	allowance, _ := strconv.ParseFloat(c.PostForm("allowance"), 64)
	bonus, _ := strconv.ParseFloat(c.PostForm("bonus"), 64)

	pf := c.PostForm("pf") != ""
	esi := c.PostForm("esi") != ""
	tds := c.PostForm("tds") != ""

	db, _ := utils.ReadDB()

	// update if exists
	for i, s := range db.SalaryStructures {
		if s.EmployeeID == employeeID {
			db.SalaryStructures[i] = models.SalaryStructure{
				EmployeeID: employeeID,
				Basic:      basic,
				HRA:        hra,
				Allowance:  allowance,
				Bonus:      bonus,
				PFEnabled:  pf,
				ESIEnabled: esi,
				TDSEnabled: tds,
				CreatedAt:  time.Now().Format("2006-01-02"),
			}

			utils.WriteDB(db)
			c.Redirect(http.StatusFound, "/admin/payroll")
			return
		}
	}

	// create new
	db.SalaryStructures = append(db.SalaryStructures, models.SalaryStructure{
		EmployeeID: employeeID,
		Basic:      basic,
		HRA:        hra,
		Allowance:  allowance,
		Bonus:      bonus,
		PFEnabled:  pf,
		ESIEnabled: esi,
		TDSEnabled: tds,
		CreatedAt:  time.Now().Format("2006-01-02"),
	})

	utils.WriteDB(db)
	c.Redirect(http.StatusFound, "/admin/payroll")
}

// ================================
// GENERATE MONTHLY PAYROLL
// ================================
func GenerateMonthlyPayroll(c *gin.Context) {

	utils.DBMutex.Lock()
	defer utils.DBMutex.Unlock()

	month := c.PostForm("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	db, _ := utils.ReadDB()

	for _, s := range db.SalaryStructures {

		// prevent duplicate payroll
		exists := false
		for _, p := range db.Payrolls {
			if p.EmployeeID == s.EmployeeID && p.Month == month {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		pay := utils.CalculatePayroll(s)
		pay.ID = len(db.Payrolls) + 1
		pay.EmployeeID = s.EmployeeID
		pay.Month = month
		pay.GeneratedAt = time.Now().Format("2006-01-02 15:04")

		db.Payrolls = append(db.Payrolls, pay)
	}

	utils.WriteDB(db)
	c.Redirect(http.StatusFound, "/admin/payroll")
}

// ================================
// HR PAYROLL LIST PAGE
// ================================
func HRPayrollPage(c *gin.Context) {
	db, _ := utils.ReadDB()

	c.HTML(http.StatusOK, "hr_payroll.html", gin.H{
		"payrolls":  db.Payrolls,
		"employees": db.Employees,
	})
}

// ================================
// VIEW PAYSLIP (HTML)
// ================================
func ViewPayslip(c *gin.Context) {
	employeeID := c.Param("employeeId")
	month := c.Param("month")

	db, _ := utils.ReadDB()

	var payroll models.Payroll
	for _, p := range db.Payrolls {
		if p.EmployeeID == employeeID && p.Month == month {
			payroll = p
			break
		}
	}

	var emp models.Employee
	for _, e := range db.Employees {
		if e.EmployeeID == employeeID {
			emp = e
			break
		}
	}

	c.HTML(http.StatusOK, "payslip_view.html", gin.H{
		"payroll":  payroll,
		"employee": emp,
	})
}

// ================================
// DOWNLOAD PAYSLIP (PDF)
// ================================
func DownloadPayslip(c *gin.Context) {
	employeeID := c.Param("employeeId")
	month := c.Param("month")

	db, _ := utils.ReadDB()

	var payroll models.Payroll
	found := false

	for _, p := range db.Payrolls {
		if p.EmployeeID == employeeID && p.Month == month {
			payroll = p
			found = true
			break
		}
	}

	if !found {
		c.String(http.StatusNotFound, "Payslip not found")
		return
	}

	var emp models.Employee
	for _, e := range db.Employees {
		if e.EmployeeID == employeeID {
			emp = e
			break
		}
	}

	path, err := utils.GeneratePayslipPDF(payroll, emp)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	c.FileAttachment(path, "Payslip_"+employeeID+"_"+month+".pdf")
}
