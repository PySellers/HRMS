package utils

import (
	"fmt"

	"pysellers-erp-go/models"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePayslipPDF(p models.Payroll, emp models.Employee) (string, error) {

	filePath := fmt.Sprintf("data/payslip_%s_%s.pdf", p.EmployeeID, p.Month)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "PySellers Digital Suite")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 8, fmt.Sprintf("Payslip for %s", p.Month))
	pdf.Ln(10)

	pdf.Cell(95, 8, "Employee Name:")
	pdf.Cell(95, 8, emp.Name)
	pdf.Ln(6)

	pdf.Cell(95, 8, "Employee ID:")
	pdf.Cell(95, 8, p.EmployeeID)
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(95, 8, "Earnings")
	pdf.Cell(95, 8, "Deductions")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 11)

	pdf.Cell(95, 8, fmt.Sprintf("Basic: %.2f", p.Basic))
	pdf.Cell(95, 8, fmt.Sprintf("PF: %.2f", p.PF))
	pdf.Ln(6)

	pdf.Cell(95, 8, fmt.Sprintf("HRA: %.2f", p.HRA))
	pdf.Cell(95, 8, fmt.Sprintf("ESI: %.2f", p.ESI))
	pdf.Ln(6)

	pdf.Cell(95, 8, fmt.Sprintf("Allowance: %.2f", p.Allowance))
	pdf.Cell(95, 8, fmt.Sprintf("TDS: %.2f", p.TDS))
	pdf.Ln(6)

	pdf.Cell(95, 8, fmt.Sprintf("Bonus: %.2f", p.Bonus))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(95, 8, fmt.Sprintf("Gross Pay: %.2f", p.Gross))
	pdf.Cell(95, 8, fmt.Sprintf("Net Pay: %.2f", p.Net))
	pdf.Ln(15)

	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(190, 8, "This is a system generated payslip.")

	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
