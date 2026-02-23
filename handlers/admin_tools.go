package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------- Show Admin Tools ----------
func ShowAdminTools(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_tools.html", gin.H{})
}

// ---------- Backup All JSON Data ----------
func BackupData(c *gin.Context) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := "backup"
	os.MkdirAll(backupDir, 0755)

	files := []string{"data/db.json", "data/finance.json", "data/recruitment.json", "data/marketing.json", "data/config.json"}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			input, _ := os.ReadFile(file)
			base := filepath.Base(file)
			os.WriteFile(filepath.Join(backupDir, timestamp+"_"+base), input, 0644)
		}
	}

	c.String(http.StatusOK, "Backup completed successfully at "+backupDir)
}

// ---------- Export Data as JSON ----------
func ExportData(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		c.String(http.StatusBadRequest, "Missing file parameter")
		return
	}
	path := filepath.Join("data", file)
	c.Header("Content-Disposition", "attachment; filename="+file)
	c.Header("Content-Type", "application/json")
	c.File(path)
}

// ---------- Clear Dummy Data ----------
func ClearDummyData(c *gin.Context) {
	files := []string{"data/db.json", "data/finance.json", "data/recruitment.json", "data/marketing.json"}
	for _, f := range files {
		_ = os.WriteFile(f, []byte("{}"), 0644)
	}
	c.String(http.StatusOK, "Dummy data cleared successfully.")
}

// ---------- System Health Check ----------
func SystemHealth(c *gin.Context) {
	status := map[string]string{
		"status":      "running",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"environment": "local development",
	}
	json.NewEncoder(c.Writer).Encode(status)
}
