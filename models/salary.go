package models

type SalaryStructure struct {
	EmployeeID string

	Basic     float64
	HRA       float64
	Allowance float64
	Bonus     float64

	PFEnabled  bool
	ESIEnabled bool
	TDSEnabled bool

	CreatedAt string
}
