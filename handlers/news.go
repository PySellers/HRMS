package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

// Show all news & announcements (Admin view)
func ShowNewsPage(c *gin.Context) {
	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	c.HTML(http.StatusOK, "news.html", gin.H{
		"news": db.News,
	})
}

// Add a new news post
func AddNews(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	news := models.News{
		ID:          len(db.News) + 1,
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Date:        time.Now().Format("2006-01-02"),
		PostedBy:    user.(string),
	}

	db.News = append(db.News, news)
	saveDB(db)

	c.Redirect(http.StatusFound, "/admin/news")
}

// Show all HR requests (Admin view)
func ShowHRRequests(c *gin.Context) {
	data, _ := os.ReadFile(dbFile)
	var db models.DB
	_ = json.Unmarshal(data, &db)

	c.HTML(http.StatusOK, "hr_requests.html", gin.H{
		"requests": db.HRRequests,
	})
}

// Update HR request status (Admin)
func UpdateHRRequestStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	action := c.Param("action") // "resolve" or "reject"

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	json.Unmarshal(data, &db)

	for i, req := range db.HRRequests {
		if req.ID == id {
			if action == "resolve" {
				db.HRRequests[i].Status = "resolved"
			} else if action == "reject" {
				db.HRRequests[i].Status = "rejected"
			}
			break
		}
	}

	saveDB(db)
	c.Redirect(http.StatusFound, "/admin/hrrequests")
}

// Employee: show HR request form
func ShowEmployeeHRForm(c *gin.Context) {
	c.HTML(http.StatusOK, "employee_hr.html", nil)
}

// Employee: submit HR request
func SubmitHRRequest(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	json.Unmarshal(data, &db)

	var employeeID int
	for _, e := range db.Employees {
		if e.Name == user.(string) {
			employeeID = e.ID
			break
		}
	}

	newRequest := models.HRRequest{
		ID:          len(db.HRRequests) + 1,
		EmployeeID:  employeeID,
		Subject:     c.PostForm("subject"),
		Description: c.PostForm("description"),
		Status:      "pending",
		Date:        time.Now().Format("2006-01-02"),
	}

	db.HRRequests = append(db.HRRequests, newRequest)
	saveDB(db)

	c.Redirect(http.StatusFound, "/employee/dashboard")
}
