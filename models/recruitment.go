package models

type JobPosting struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Department     string `json:"department"`
	Location       string `json:"location"`
	EmploymentType string `json:"employment_type"` // Full-time, Part-time, Contract, Internship
	Experience     string `json:"experience"`      // e.g., "2-4 years"
	SalaryRange    string `json:"salary_range"`    // e.g., "₹4L–₹6L per annum"
	PostedBy       string `json:"posted_by"`
	PostedDate     string `json:"posted_date"`
	Description    string `json:"description"`
	Status         string `json:"status"`          // Open / Closed
}
