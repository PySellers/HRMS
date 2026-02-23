package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

var recruitmentFile = "data/recruitment.json"

// ---------- Show All Jobs ----------
func ShowRecruitments(c *gin.Context) {
	data, _ := os.ReadFile(recruitmentFile)
	var jobs []models.JobPosting
	_ = json.Unmarshal(data, &jobs)

	c.HTML(http.StatusOK, "recruitment.html", gin.H{
		"jobs": jobs,
	})
}

// ---------- Add New Job ----------
func AddRecruitment(c *gin.Context) {
	data, _ := os.ReadFile(recruitmentFile)
	var jobs []models.JobPosting
	_ = json.Unmarshal(data, &jobs)

	newJob := models.JobPosting{
		ID:             len(jobs) + 1,
		Title:          c.PostForm("title"),
		Department:     c.PostForm("department"),
		Location:       c.PostForm("location"),
		EmploymentType: c.PostForm("employment_type"),
		Experience:     c.PostForm("experience"),
		SalaryRange:    c.PostForm("salary_range"),
		PostedBy:       "Admin",
		PostedDate:     time.Now().Format("2006-01-02"),
		Description:    c.PostForm("description"),
		Status:         "Open",
	}

	jobs = append(jobs, newJob)
	saveRecruitment(jobs)

	c.Redirect(http.StatusFound, "/admin/recruitment")
}

// ---------- Update Job Status ----------
func UpdateRecruitmentStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	action := c.Param("action") // "close" or "open"

	data, _ := os.ReadFile(recruitmentFile)
	var jobs []models.JobPosting
	_ = json.Unmarshal(data, &jobs)

	for i, job := range jobs {
		if job.ID == id {
			if action == "close" {
				jobs[i].Status = "Closed"
			} else if action == "open" {
				jobs[i].Status = "Open"
			}
		}
	}

	saveRecruitment(jobs)
	c.Redirect(http.StatusFound, "/admin/recruitment")
}

// ---------- Helper ----------
func saveRecruitment(jobs []models.JobPosting) {
	data, _ := json.MarshalIndent(jobs, "", "  ")
	os.WriteFile(recruitmentFile, data, 0644)
}
