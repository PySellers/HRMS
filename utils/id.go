package utils

import "pysellers-erp-go/models"

func NextEmployeeID(db *models.DB) int {
	max := 0
	for _, e := range db.Employees {
		if e.ID > max {
			max = e.ID
		}
	}
	return max + 1
}

func NextUserID(db *models.DB) int {
	max := 0
	for _, u := range db.Users {
		if u.ID > max {
			max = u.ID
		}
	}
	return max + 1
}
