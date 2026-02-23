package models

// Config stores company & system settings
type Config struct {
    CompanyName       string `json:"company_name"`
    Logo              string `json:"logo"`
    Address           string `json:"address"`
    GSTNumber         string `json:"gst_number"`
    DefaultLeaveLimit int    `json:"default_leave_limit"`
    WorkingHours      string `json:"working_hours"`
    BackupSchedule    string `json:"backup_schedule"`
    TimeZone          string `json:"time_zone"`
}
