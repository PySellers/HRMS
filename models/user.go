package models

type User struct {
	ID                  int    `json:"id"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	Role                string `json:"role"` // "admin" or "employee"
	ForceChangePassword bool   `json:"force_change_password"`
}
