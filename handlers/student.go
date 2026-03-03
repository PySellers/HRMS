package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"pysellers-erp-go/models"

	"github.com/gin-gonic/gin"
)

// ShowStudentTraining - displays student training workspace
func ShowStudentTraining(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Debug print
	fmt.Printf("Student login - User ID: %d\n", user.ID)

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error":   "Failed to load training data",
			"details": err.Error(),
		})
		return
	}

	mainDB, err := loadMainDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error":   "Failed to load main data",
			"details": err.Error(),
		})
		return
	}

	// Debug: Print all students
	fmt.Println("Students in DB:")
	for _, s := range trainingDB.Students {
		fmt.Printf("  Student ID: %d, Name: %s, BatchID: %d\n", s.ID, s.Name, s.BatchID)
	}

	// Get student details
	var student models.Student
	var studentFound bool
	for _, s := range trainingDB.Students {
		if s.ID == user.ID {
			student = s
			studentFound = true
			fmt.Printf("Found student: %s (ID: %d, Batch: %d)\n", student.Name, student.ID, student.BatchID)
			break
		}
	}

	if !studentFound {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error":   "Student record not found",
			"details": fmt.Sprintf("User ID %d not found in training database. Please contact your mentor.", user.ID),
		})
		return
	}

	// Get student's batch details
	var studentBatch models.Batch
	var mentorName string
	var mentorID int
	batchFound := false
	for _, batch := range trainingDB.Batches {
		if batch.ID == student.BatchID {
			studentBatch = batch
			mentorID = batch.MentorID
			batchFound = true
			// Get mentor name from employees
			for _, emp := range mainDB.Employees {
				if emp.ID == batch.MentorID {
					mentorName = emp.Name
					break
				}
			}
			break
		}
	}

	if !batchFound {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error":   "Batch not found",
			"details": fmt.Sprintf("Batch ID %d not found for student %s", student.BatchID, student.Name),
		})
		return
	}

	// Get student name from employees
	studentName := student.Name
	for _, emp := range mainDB.Employees {
		if emp.ID == user.ID {
			studentName = emp.Name
			break
		}
	}

	// Get all sessions for student's batch
	var completedSessions []map[string]interface{}
	var upcomingSessions []map[string]interface{}
	var allSessions []map[string]interface{}
	var totalSessions, attendedSessions int

	now := time.Now()
	for _, session := range trainingDB.Sessions {
		if session.BatchID == student.BatchID {
			sessionDate, _ := time.Parse("2006-01-02", session.Date)

			present := false
			if val, ok := session.Attendance[user.ID]; ok {
				present = val
			}

			// Get trainer name
			trainerName := session.Trainer
			trainerID, _ := strconv.Atoi(session.Trainer)
			for _, emp := range mainDB.Employees {
				if emp.ID == trainerID {
					trainerName = emp.Name
					break
				}
			}

			sessionData := map[string]interface{}{
				"ID":      session.ID,
				"Date":    session.Date,
				"Topic":   session.Topic,
				"Trainer": trainerName,
				"Notes":   session.Notes,
				"Status":  session.Status,
				"Present": present,
			}

			allSessions = append(allSessions, sessionData)

			if sessionDate.Before(now) || sessionDate.Equal(now) {
				completedSessions = append(completedSessions, sessionData)
				totalSessions++
				if present {
					attendedSessions++
				}
			} else {
				upcomingSessions = append(upcomingSessions, sessionData)
			}
		}
	}

	// Calculate attendance percentage
	attendancePercent := 0.0
	if totalSessions > 0 {
		attendancePercent = float64(attendedSessions) / float64(totalSessions) * 100
	}

	// Get all materials for student's batch
	var courseMaterials []map[string]interface{}
	for _, material := range trainingDB.Materials {
		if material.BatchID == student.BatchID {
			courseMaterials = append(courseMaterials, map[string]interface{}{
				"ID":       material.ID,
				"Title":    material.Title,
				"Type":     material.Type,
				"URL":      material.URL,
				"Uploaded": material.Uploaded,
			})
		}
	}

	// Get assignments for student
	var studentAssignments []map[string]interface{}
	var totalAssignments, completedAssignments, gradedAssignments int
	var totalMarks int

	for _, assignment := range trainingDB.Assignments {
		// Check if assignment is for student's batch or specifically for this student
		if assignment.BatchID == student.BatchID || assignment.StudentID == user.ID {
			totalAssignments++

			// Check if submitted
			var submission *models.Submission
			var grade *models.Grade

			if trainingDB.Submissions != nil {
				for _, sub := range trainingDB.Submissions {
					if sub.AssignmentID == assignment.ID && sub.StudentID == user.ID {
						submission = &sub
						break
					}
				}
			}

			// Get grade if exists
			if submission != nil && trainingDB.Grades != nil {
				for _, g := range trainingDB.Grades {
					if g.AssignmentID == assignment.ID && g.StudentID == user.ID {
						grade = &g
						gradedAssignments++
						totalMarks += g.Score
						break
					}
				}
				completedAssignments++
			} else if submission != nil {
				completedAssignments++
			}

			// Check if deadline passed
			dueDate, _ := time.Parse("2006-01-02", assignment.DueDate)
			isDeadlinePassed := dueDate.Before(now)

			// Get created by (mentor) name
			createdByName := ""
			for _, emp := range mainDB.Employees {
				if emp.ID == assignment.CreatedBy {
					createdByName = emp.Name
					break
				}
			}

			studentAssignments = append(studentAssignments, map[string]interface{}{
				"ID":            assignment.ID,
				"Title":         assignment.Title,
				"Description":   assignment.Description,
				"DueDate":       assignment.DueDate,
				"DueDatePassed": isDeadlinePassed,
				"Submitted":     submission != nil,
				"Submission":    submission,
				"Grade":         grade,
				"MaxScore":      assignment.MaxScore,
				"CreatedBy":     createdByName,
				"CreatedAt":     assignment.CreatedAt,
			})
		}
	}

	// Calculate average marks
	averageMarks := 0.0
	if gradedAssignments > 0 {
		averageMarks = float64(totalMarks) / float64(gradedAssignments)
	}

	// Calculate course progress
	courseProgress := 0.0
	if totalSessions > 0 || totalAssignments > 0 {
		sessionWeight := 0.5
		assignmentWeight := 0.5

		sessionProgress := 0.0
		if totalSessions > 0 {
			sessionProgress = float64(attendedSessions) / float64(totalSessions) * 100 * sessionWeight
		}

		assignmentProgress := 0.0
		if totalAssignments > 0 {
			assignmentProgress = float64(completedAssignments) / float64(totalAssignments) * 100 * assignmentWeight
		}

		courseProgress = sessionProgress + assignmentProgress
	}

	// Get suggested next course based on progress
	nextCourse := ""
	if courseProgress >= 80 {
		pathways := map[string]string{
			"Python Full Stack": "Advanced Python with Django",
			"Java Full Stack":   "Advanced Java with Spring Boot",
			"Data Science":      "Machine Learning with Python",
			"AI":                "Deep Learning & Neural Networks",
		}

		if next, exists := pathways[student.Course]; exists {
			nextCourse = next
		} else {
			nextCourse = "Advanced " + student.Course
		}
	}

	// Mock fee details
	feeDetails := map[string]interface{}{
		"TotalFee":      50000,
		"PaidAmount":    25000,
		"PendingAmount": 25000,
		"Status":        "Partial",
		"DueDate":       "2026-03-15",
		"Transactions": []map[string]interface{}{
			{
				"Date":   "2026-01-15",
				"Amount": 25000,
				"Mode":   "Online Transfer",
			},
		},
	}

	// Check certificate eligibility
	certificateEligible := attendancePercent >= 80 && completedAssignments == totalAssignments && totalAssignments > 0

	// Get any certificates
	var studentCertificate map[string]interface{}
	if trainingDB.Certificates != nil {
		for _, cert := range trainingDB.Certificates {
			if cert.StudentID == user.ID {
				studentCertificate = map[string]interface{}{
					"ID":             cert.ID,
					"IssueDate":      cert.IssueDate,
					"CertificateURL": cert.CertificateURL,
					"IsIssued":       cert.IsIssued,
				}
				break
			}
		}
	}

	// Calculate pending assignments
	pendingAssignments := totalAssignments - completedAssignments

	c.HTML(http.StatusOK, "training_student.html", gin.H{
		"title":                "Student Learning Dashboard",
		"Student":              gin.H{"Name": studentName, "ID": user.ID},
		"Course":               student.Course,
		"Batch":                studentBatch,
		"MentorName":           mentorName,
		"MentorID":             mentorID,
		"Attendance":           attendancePercent,
		"AttendedSessions":     attendedSessions,
		"TotalSessions":        totalSessions,
		"CompletedSessions":    completedSessions,
		"UpcomingSessions":     upcomingSessions,
		"AllSessions":          allSessions,
		"Assignments":          studentAssignments,
		"TotalAssignments":     totalAssignments,
		"CompletedAssignments": completedAssignments,
		"PendingAssignments":   pendingAssignments,
		"GradedAssignments":    gradedAssignments,
		"AverageMarks":         averageMarks,
		"CourseProgress":       courseProgress,
		"NextCourse":           nextCourse,
		"CourseMaterials":      courseMaterials,
		"FeeDetails":           feeDetails,
		"Certificate":          studentCertificate,
		"CertificateEligible":  certificateEligible,
		"Today":                now.Format("2006-01-02"),
	})
}

// Upload Assignment Submission
func UploadSubmission(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	assignmentID, _ := strconv.Atoi(c.PostForm("assignment_id"))

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	// Check if assignment exists and deadline not passed
	var assignment models.Assignment
	for _, a := range trainingDB.Assignments {
		if a.ID == assignmentID {
			assignment = a
			break
		}
	}

	if assignment.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Check deadline
	dueDate, _ := time.Parse("2006-01-02", assignment.DueDate)
	if time.Now().After(dueDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Submission deadline has passed"})
		return
	}

	// Check if already submitted
	if trainingDB.Submissions != nil {
		for _, sub := range trainingDB.Submissions {
			if sub.AssignmentID == assignmentID && sub.StudentID == user.ID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Already submitted. No editing allowed."})
				return
			}
		}
	}

	// Handle file upload
	file, err := c.FormFile("submission_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Create submissions directory if not exists
	uploadDir := "uploads/submissions"
	os.MkdirAll(uploadDir, 0755)

	// Save file with timestamp to prevent duplicates
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Base(file.Filename)
	safeFilename := fmt.Sprintf("%d_%s_%s", user.ID, timestamp, filename)
	filepath := filepath.Join(uploadDir, safeFilename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create submission record
	newSubmission := models.Submission{
		ID:             len(trainingDB.Submissions) + 1,
		AssignmentID:   assignmentID,
		StudentID:      user.ID,
		SubmittedAt:    time.Now().Format("2006-01-02 15:04:05"),
		Files:          []string{filepath},
		Status:         "submitted",
		LateSubmission: false,
	}

	if trainingDB.Submissions == nil {
		trainingDB.Submissions = []models.Submission{newSubmission}
	} else {
		trainingDB.Submissions = append(trainingDB.Submissions, newSubmission)
	}

	saveTrainingDB(trainingDB)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Assignment submitted successfully",
	})
}

// Get Assignment Details
func GetAssignmentDetails(c *gin.Context) {
	assignmentID, _ := strconv.Atoi(c.Param("id"))
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	var assignment models.Assignment
	var submission *models.Submission
	var grade *models.Grade

	// Find assignment
	for _, a := range trainingDB.Assignments {
		if a.ID == assignmentID {
			assignment = a
			break
		}
	}

	if assignment.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Find submission
	if trainingDB.Submissions != nil {
		for _, sub := range trainingDB.Submissions {
			if sub.AssignmentID == assignmentID && sub.StudentID == user.ID {
				submission = &sub
				break
			}
		}
	}

	// Find grade
	if submission != nil && trainingDB.Grades != nil {
		for _, g := range trainingDB.Grades {
			if g.AssignmentID == assignmentID && g.StudentID == user.ID {
				grade = &g
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"assignment": assignment,
		"submission": submission,
		"grade":      grade,
	})
}

// Get Session Details
func GetSessionDetails(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	var session models.Session
	for _, s := range trainingDB.Sessions {
		if s.ID == sessionID {
			session = s
			break
		}
	}

	if session.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	present := false
	if val, ok := session.Attendance[user.ID]; ok {
		present = val
	}

	c.JSON(http.StatusOK, gin.H{
		"session": session,
		"present": present,
	})
}

// Request Certificate
func RequestCertificate(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	// Check if already requested
	if trainingDB.Certificates != nil {
		for _, cert := range trainingDB.Certificates {
			if cert.StudentID == user.ID {
				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate already requested"})
				return
			}
		}
	}

	// Get student course
	course := ""
	for _, student := range trainingDB.Students {
		if student.ID == user.ID {
			course = student.Course
			break
		}
	}

	// Create certificate request
	newCert := models.Certificate{
		ID:             len(trainingDB.Certificates) + 1,
		StudentID:      user.ID,
		Course:         course,
		IssueDate:      "",
		CertificateURL: "",
		IsEligible:     true,
		IsIssued:       false,
	}

	if trainingDB.Certificates == nil {
		trainingDB.Certificates = []models.Certificate{newCert}
	} else {
		trainingDB.Certificates = append(trainingDB.Certificates, newCert)
	}

	saveTrainingDB(trainingDB)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Certificate request submitted successfully",
	})
}
