package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"pysellers-erp-go/models"
)

const dbPath = "data/db.json"

func main() {
	data, err := os.ReadFile(dbPath)
	if err != nil {
		fmt.Println("failed to read db:", err)
		os.Exit(1)
	}

	var db models.DB
	if err := json.Unmarshal(data, &db); err != nil {
		fmt.Println("failed to parse db:", err)
		os.Exit(1)
	}

	// Backup original
	ts := time.Now().Format("20060102_150405")
	backupDir := "backup"
	_ = os.MkdirAll(backupDir, 0755)
	backupPath := filepath.Join(backupDir, ts+"_db.json")
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		fmt.Println("warning: failed to write backup:", err)
	} else {
		fmt.Println("backup written to", backupPath)
	}

	// Deduplicate employees by EmployeeID and Email (keep first occurrence)
	seenEmpID := make(map[string]bool)
	seenEmail := make(map[string]bool)
	newEmps := make([]models.Employee, 0, len(db.Employees))
	for _, e := range db.Employees {
		id := normalize(e.EmployeeID)
		em := normalize(e.Email)
		if id != "" && seenEmpID[id] {
			continue
		}
		if em != "" && seenEmail[em] {
			continue
		}
		if id != "" {
			seenEmpID[id] = true
		}
		if em != "" {
			seenEmail[em] = true
		}
		newEmps = append(newEmps, e)
	}

	// Deduplicate users by username (keep first occurrence)
	seenUser := make(map[string]bool)
	newUsers := make([]models.User, 0, len(db.Users))
	for _, u := range db.Users {
		uname := normalize(u.Username)
		if uname == "" {
			// keep empty username entries as-is
			newUsers = append(newUsers, u)
			continue
		}
		if seenUser[uname] {
			continue
		}
		seenUser[uname] = true
		newUsers = append(newUsers, u)
	}

	db.Employees = newEmps
	db.Users = newUsers

	out, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal cleaned db:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(dbPath, out, 0644); err != nil {
		fmt.Println("failed to write cleaned db:", err)
		os.Exit(1)
	}

	fmt.Println("dedupe complete. employees:", len(newEmps), "users:", len(newUsers))
}

func normalize(s string) string {
	// simple lowercase + trim
	return stringsTrimSpaceToLower(s)
}

// small inline helper to avoid importing strings package twice
func stringsTrimSpaceToLower(s string) string {
	// replicate strings.TrimSpace + strings.ToLower
	// simple implementation:
	b := []byte(s)
	// trim left
	start := 0
	for start < len(b) && (b[start] == ' ' || b[start] == '\t' || b[start] == '\n' || b[start] == '\r') {
		start++
	}
	end := len(b)
	for end > start && (b[end-1] == ' ' || b[end-1] == '\t' || b[end-1] == '\n' || b[end-1] == '\r') {
		end--
	}
	if start >= end {
		return ""
	}
	out := make([]byte, end-start)
	for i := start; i < end; i++ {
		ch := b[i]
		if ch >= 'A' && ch <= 'Z' {
			ch = ch + ('a' - 'A')
		}
		out[i-start] = ch
	}
	return string(out)
}
