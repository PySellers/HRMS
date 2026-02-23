package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

var trainingFile = "data/training.json"

// ---------- Show Dashboard ----------
func ShowTrainingDashboard(c *gin.Context) {
	data, _ := os.ReadFile(trainingFile)
	var db models.TrainingDB
	json.Unmarshal(data, &db)

	c.HTML(http.StatusOK, "training.html", gin.H{
		"students": db.Students,
		"batches":  db.Batches,
		"sessions": db.Sessions,
	})
}

// ---------- Add Student ----------
func AddStudent(c *gin.Context) {
	data, _ := os.ReadFile(trainingFile)
	var db models.TrainingDB
	json.Unmarshal(data, &db)

	newStudent := models.Student{
		ID:       len(db.Students) + 1,
		Name:     c.PostForm("name"),
		Email:    c.PostForm("email"),
		Phone:    c.PostForm("phone"),
		Course:   c.PostForm("course"),
		BatchID:  len(db.Batches), // assign to latest batch
		JoinDate: time.Now().Format("2006-01-02"),
		Status:   "Active",
	}
	db.Students = append(db.Students, newStudent)
	saveTrainingDB(db)

	c.Redirect(http.StatusFound, "/admin/training")
}

// ---------- Add Session ----------
func AddSession(c *gin.Context) {
	data, _ := os.ReadFile(trainingFile)
	var db models.TrainingDB
	json.Unmarshal(data, &db)

	newSession := models.Session{
		ID:        len(db.Sessions) + 1,
		BatchID:   len(db.Batches),
		Date:      time.Now().Format("2006-01-02"),
		Topic:     c.PostForm("topic"),
		Trainer:   c.PostForm("trainer"),
		Attendance: make(map[int]bool),
	}
	db.Sessions = append(db.Sessions, newSession)
	saveTrainingDB(db)

	c.Redirect(http.StatusFound, "/admin/training")
}

// ---------- Add Material ----------
func AddMaterial(c *gin.Context) {
	data, _ := os.ReadFile(trainingFile)
	var db models.TrainingDB
	json.Unmarshal(data, &db)

	newMaterial := models.Material{
		ID:       len(db.Materials) + 1,
		BatchID:  len(db.Batches),
		Title:    c.PostForm("title"),
		URL:      c.PostForm("url"),
		Type:     c.PostForm("type"),
		Uploaded: time.Now().Format("2006-01-02"),
	}
	db.Materials = append(db.Materials, newMaterial)
	saveTrainingDB(db)

	c.Redirect(http.StatusFound, "/admin/training")
}

// ---------- Save ----------
func saveTrainingDB(db models.TrainingDB) {
	data, _ := json.MarshalIndent(db, "", "  ")
	_ = os.WriteFile(trainingFile, data, 0644)
}

// ShowMentorTraining - displays mentor’s training management page
func ShowMentorTraining(c *gin.Context) {
    c.HTML(http.StatusOK, "training_mentor.html", gin.H{
        "title": "Mentor Training Dashboard",
    })
}

// ShowStudentTraining - displays student training workspace
func ShowStudentTraining(c *gin.Context) {
    c.HTML(http.StatusOK, "training_student.html", gin.H{
        "title": "Student Learning Dashboard",
    })
}
