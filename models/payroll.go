package models

type Payroll struct {
    ID         int     `json:"id"`
    EmployeeID int     `json:"employee_id"`
    Month      string  `json:"month"`   // YYYY-MM
    NetSalary  float64 `json:"net_salary"`
    Status     string  `json:"status"`  // paid, pending
}
