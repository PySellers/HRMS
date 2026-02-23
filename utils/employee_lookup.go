package utils

import "pysellers-erp-go/models"

// Username == email OR employee_id (your system supports both)
func GetEmployeeByUsername(db *models.DB, username string) *models.Employee {
	for _, e := range db.Employees {
		if e.Email == username || e.EmployeeID == username {
			return &e
		}
	}
	return nil
}
