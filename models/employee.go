package models

type Project struct {
	Name      string   `json:"name"`
	Manager   string   `json:"manager"`
	Coworkers []string `json:"coworkers"`
}

type Employee struct {
	ID             int      `json:"id"`
	EmployeeID     string   `json:"employee_id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	Role           string   `json:"role"`
	Address        string   `json:"address"`
	IsActive       bool     `json:"is_active"`
	ProfilePic     string   `json:"profile_pic"`
	Project        Project  `json:"project"`
	Skills         []string `json:"skills"`
	Certifications []string `json:"certifications"`
}
