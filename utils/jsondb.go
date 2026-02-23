package utils

import (
	"encoding/json"
	"os"

	"pysellers-erp-go/models"
)

var DBFile = "data/db.json"

func ReadDB() (*models.DB, error) {
	data, err := os.ReadFile(DBFile)
	if err != nil {
		return nil, err
	}

	var db models.DB
	err = json.Unmarshal(data, &db)
	return &db, err
}

func WriteDB(db *models.DB) error {
	out, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(DBFile, out, 0644)
}
