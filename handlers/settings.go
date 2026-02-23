package handlers

import (
	"encoding/json"
	"net/http"
	"os"
    "strconv"
	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

var configFile = "data/config.json"

// ---------- Show Settings ----------
func ShowSettings(c *gin.Context) {
	data, _ := os.ReadFile(configFile)
	var cfg models.Config
	_ = json.Unmarshal(data, &cfg)

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"config": cfg,
	})
}

// ---------- Update Settings ----------
func UpdateSettings(c *gin.Context) {
	var cfg models.Config

	cfg.CompanyName = c.PostForm("company_name")
	cfg.Logo = c.PostForm("logo")
	cfg.Address = c.PostForm("address")
	cfg.GSTNumber = c.PostForm("gst_number")
	cfg.DefaultLeaveLimit, _ = strconv.Atoi(c.PostForm("default_leave_limit"))
	cfg.WorkingHours = c.PostForm("working_hours")
	cfg.BackupSchedule = c.PostForm("backup_schedule")
	cfg.TimeZone = c.PostForm("time_zone")

	data, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(configFile, data, 0644)

	c.Redirect(http.StatusFound, "/admin/settings")
}
