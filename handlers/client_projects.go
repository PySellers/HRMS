package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"pysellers-erp-go/models"

	"github.com/gin-gonic/gin"
)

var clientDBFile = "data/client_projects.json"

// Show Client & Project Dashboard
func ShowClientDashboard(c *gin.Context) {
	data, _ := os.ReadFile(clientDBFile)
	var db models.ClientDB
	json.Unmarshal(data, &db)

	c.HTML(http.StatusOK, "client_projects.html", gin.H{
		"clients":  db.Clients,
		"projects": db.Projects,
		"tasks":    db.Tasks,
		"invoices": db.Invoices,
	})
}

// Add Client
func AddClient(c *gin.Context) {
	data, _ := os.ReadFile(clientDBFile)
	var db models.ClientDB
	json.Unmarshal(data, &db)

	newClient := models.Client{
		ID:      len(db.Clients) + 1,
		Name:    c.PostForm("name"),
		Email:   c.PostForm("email"),
		Phone:   c.PostForm("phone"),
		Company: c.PostForm("company"),
		Address: c.PostForm("address"),
		Notes:   c.PostForm("notes"),
	}
	db.Clients = append(db.Clients, newClient)

	saveClientDB(db)
	c.Redirect(http.StatusFound, "/admin/clients")
}

// Add Project
func AddProject(c *gin.Context) {
	data, _ := os.ReadFile(clientDBFile)
	var db models.ClientDB
	json.Unmarshal(data, &db)

	clientID, _ := strconv.Atoi(c.PostForm("client_id"))
	budget, _ := strconv.ParseFloat(c.PostForm("budget"), 64)

	newProject := models.ClientProject{
		ID:          len(db.Projects) + 1,
		ClientID:    clientID,
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Status:      "Planned",
		StartDate:   c.PostForm("start_date"),
		EndDate:     c.PostForm("end_date"),
		Budget:      budget,
	}
	db.Projects = append(db.Projects, newProject)

	saveClientDB(db)
	c.Redirect(http.StatusFound, "/admin/clients")
}

// Add Task (for Kanban board)
func AddTask(c *gin.Context) {
	data, _ := os.ReadFile(clientDBFile)
	var db models.ClientDB
	json.Unmarshal(data, &db)

	projectID, _ := strconv.Atoi(c.PostForm("project_id"))

	newTask := models.Task{
		ID:          len(db.Tasks) + 1,
		ProjectID:   projectID,
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Assignee:    c.PostForm("assignee"),
		Status:      "To Do",
		Priority:    c.PostForm("priority"),
		DueDate:     c.PostForm("due_date"),
	}
	db.Tasks = append(db.Tasks, newTask)

	saveClientDB(db)
	c.Redirect(http.StatusFound, "/admin/clients")
}

// Update Task Status (drag-drop simulation)
func UpdateTaskStatus(c *gin.Context) {
	data, _ := os.ReadFile(clientDBFile)
	var db models.ClientDB
	json.Unmarshal(data, &db)

	id, _ := strconv.Atoi(c.Param("id"))
	newStatus := c.PostForm("status")

	for i, t := range db.Tasks {
		if t.ID == id {
			db.Tasks[i].Status = newStatus
			break
		}
	}
	saveClientDB(db)
	c.Redirect(http.StatusFound, "/admin/clients")
}

// Generate Invoice
func GenerateInvoice(c *gin.Context) {
	data, _ := os.ReadFile(clientDBFile)
	var db models.ClientDB
	json.Unmarshal(data, &db)

	clientID, _ := strconv.Atoi(c.PostForm("client_id"))
	projectID, _ := strconv.Atoi(c.PostForm("project_id"))
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)

	newInvoice := models.Invoice{
		ID:        len(db.Invoices) + 1,
		ClientID:  clientID,
		ProjectID: projectID,
		Amount:    amount,
		Date:      time.Now().Format("2006-01-02"),
		Status:    "Pending",
	}
	db.Invoices = append(db.Invoices, newInvoice)

	saveClientDB(db)
	c.Redirect(http.StatusFound, "/admin/clients")
}

// Helper
func saveClientDB(db models.ClientDB) {
	out, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(clientDBFile, out, 0644)
}
