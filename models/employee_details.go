package models

// HR-ONLY EXTENDED EMPLOYEE DETAILS
type EmployeeDetails struct {
	EmployeeID string `json:"employee_id"`

	// =====================
	// PERSONAL DETAILS
	// =====================
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`

	// =====================
	// EMPLOYMENT DETAILS
	// =====================
	Designation    string `json:"designation"`
	Department     string `json:"department"`
	DateOfJoining  string `json:"date_of_joining"`
	EmploymentType string `json:"employment_type"`

	// =====================
	// GOVERNMENT IDS
	// =====================
	PAN     string `json:"pan"`
	Aadhaar string `json:"aadhaar"`

	// =====================
	// BANK DETAILS
	// =====================
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	IFSC          string `json:"ifsc"`

	// =====================
	// EDUCATION
	// =====================
	Degree     string `json:"degree"`
	Experience string `json:"experience"`
	Resume     string `json:"resume"`

	// =====================
	// SYSTEM
	// =====================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
