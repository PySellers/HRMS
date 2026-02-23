package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

var marketingFile = "data/marketing.json"

// ShowMarketingPage renders the form + summary
func ShowMarketingPage(c *gin.Context) {
	data, _ := os.ReadFile(marketingFile)
	var db models.MarketingDB
	json.Unmarshal(data, &db)

	// Generate simple daily/weekly summaries
	today := time.Now().Format("2006-01-02")
	dailyCount := 0
	totalValue := 0.0
	for _, a := range db.Activities {
		if a.Date == today {
			dailyCount++
			totalValue += a.EstimatedValue
		}
	}

	c.HTML(http.StatusOK, "marketing.html", gin.H{
		"activities":  db.Activities,
		"todayCount":  dailyCount,
		"totalValue":  totalValue,
	})
}

// AddMarketingActivity handles new form submissions
func AddMarketingActivity(c *gin.Context) {
	os.MkdirAll("data", os.ModePerm)

	file, _ := os.ReadFile(marketingFile)
	var db models.MarketingDB
	json.Unmarshal(file, &db)

	estimatedVal, _ := strconv.ParseFloat(c.PostForm("estimated_value"), 64)
	timeSpent, _ := strconv.Atoi(c.PostForm("time_spent"))

	newAct := models.MarketingActivity{
		ID:               len(db.Activities) + 1,
		Date:             time.Now().Format("2006-01-02"),
		PersonName:       c.PostForm("person_name"),
		Role:             c.PostForm("role"),
		OutreachSource:   strings.Join(c.PostFormArray("outreach_source"), ", "),
		PlatformName:     c.PostForm("platform_name"),
		ActivityType:     c.PostForm("activity_type"),
		Industry:         c.PostForm("industry"),
		ProjectType:      c.PostForm("project_type"),
		EstimatedValue:   estimatedVal,
		URL:              c.PostForm("url"),
		Status:           c.PostForm("status"),
		NextFollowUpDate: c.PostForm("next_followup_date"),
		KeyNotes:         c.PostForm("key_notes"),
		TimeSpent:        timeSpent,
	}

	db.Activities = append(db.Activities, newAct)

	out, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(marketingFile, out, 0644)

	c.Redirect(http.StatusFound, "/admin/marketing")
}

// GenerateWeeklySummary endpoint (for dashboard use)
func GenerateWeeklySummary(c *gin.Context) {
	data, _ := os.ReadFile(marketingFile)
	var db models.MarketingDB
	json.Unmarshal(data, &db)

	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))

	count := 0
	conversions := 0
	totalValue := 0.0

	for _, a := range db.Activities {
		actDate, _ := time.Parse("2006-01-02", a.Date)
		if actDate.After(startOfWeek) && actDate.Before(now.Add(24*time.Hour)) {
			count++
			if a.Status == "Converted" {
				conversions++
			}
			totalValue += a.EstimatedValue
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"weekly_total":     count,
		"converted":        conversions,
		"conversion_ratio": float64(conversions) / float64(count+1),
		"pipeline_value":   totalValue,
	})
}
