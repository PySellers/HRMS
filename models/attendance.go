package models

type Attendance struct {
	ID         int                 `json:"id"`
	EmployeeID int                 `json:"employee_id"`
	Date       string              `json:"date"`
	Sessions   []AttendanceSession `json:"sessions"`
	TotalTime  string              `json:"total_time"`
}

type AttendanceSession struct {
	TimeIn    string `json:"time_in"`
	TimeOut   string `json:"time_out,omitempty"`
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
}
