package models

type DB struct {
	Users        []User        `json:"users"`
	Employees    []Employee    `json:"employees"`
	Attendance   []Attendance  `json:"attendance"`
	Leaves       []Leave       `json:"leaves"`
	Payroll      []Payroll     `json:"payroll"`
	Transactions []Transaction `json:"transactions"`
	Finance      []Finance     `json:"finance"`
	News         []News        `json:"news"`
	HRRequests   []HRRequest   `json:"hr_requests"`
	Policies     []Policy      `json:"policies"`
}
