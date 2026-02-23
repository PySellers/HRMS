package utils

import (
	"strings"

	"pysellers-erp-go/models"
)

func EmployeeOrUserExists(db *models.DB, empID, email string) bool {
	empID = strings.ToUpper(strings.TrimSpace(empID))
	email = strings.ToLower(strings.TrimSpace(email))

	for _, e := range db.Employees {
		if strings.ToUpper(e.EmployeeID) == empID {
			return true
		}
		if strings.ToLower(e.Email) == email {
			return true
		}
	}

	for _, u := range db.Users {
		if strings.ToUpper(u.Username) == empID {
			return true
		}
	}

	return false
}
