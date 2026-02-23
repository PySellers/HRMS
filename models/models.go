// models/models.go
package models

// Leave represents an employee leave request
type Leave struct {
	ID         int    `json:"id"`
	EmployeeID int    `json:"employee_id"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	Reason     string `json:"reason"`
	Status     string `json:"status"` // pending / approved / rejected
}

// Finance represents an income or expense record
type Finance struct {
	ID          int     `json:"id"`
	Type        string  `json:"type"`        // "income" or "expense"
	Category    string  `json:"category"`    // e.g., "Salary", "Rent", etc.
	Description string  `json:"description"` // short text
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"` // YYYY-MM-DD
}

// News or announcement posted by Admin
type News struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	PostedBy    string `json:"posted_by"`
}

// HR request raised by employee
type HRRequest struct {
	ID          int    `json:"id"`
	EmployeeID  int    `json:"employee_id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Status      string `json:"status"` // pending / resolved / rejected
	Date        string `json:"date"`
}

// Policy document (uploaded or referenced)
type Policy struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
