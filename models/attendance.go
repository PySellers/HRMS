package models

type Attendance struct {
	ID         int    `json:"id"`
	EmployeeID int    `json:"employee_id"`
	Date       string `json:"date"`
	TimeIn     string `json:"time_in"`
	TimeOut    string `json:"time_out"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
}
