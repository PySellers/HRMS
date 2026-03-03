package models

type Payroll struct {
	ID         int
	EmployeeID string
	Month      string // e.g. "2026-02"

	Basic     float64
	HRA       float64
	Allowance float64
	Bonus     float64

	Gross float64

	PF  float64
	ESI float64
	TDS float64

	Net float64

	GeneratedAt string
}
