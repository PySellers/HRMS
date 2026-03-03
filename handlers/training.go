package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"pysellers-erp-go/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var trainingFile = "data/training.json"
var mainDBFile = "data/db.json"

// ---------- Helper Functions ----------
func loadTrainingDB() (models.TrainingDB, error) {
	var db models.TrainingDB
	data, err := os.ReadFile(trainingFile)
	if err != nil {
		// If file doesn't exist, create initial structure
		db = models.TrainingDB{
			Students:    []models.Student{},
			Batches:     []models.Batch{},
			Sessions:    []models.Session{},
			Materials:   []models.Material{},
			Assignments: []models.Assignment{},
			Submissions: []models.Submission{},
			Grades:      []models.Grade{},
		}
		return db, nil
	}
	err = json.Unmarshal(data, &db)
	return db, err
}

func saveTrainingDB(db models.TrainingDB) error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(trainingFile, data, 0644)
}

func getCurrentUser(c *gin.Context) (models.User, error) {
	session := sessions.Default(c)
	username := session.Get("user")

	var mainDB struct {
		Users []models.User `json:"users"`
	}

	data, err := os.ReadFile(mainDBFile)
	if err != nil {
		return models.User{}, err
	}

	err = json.Unmarshal(data, &mainDB)
	if err != nil {
		return models.User{}, err
	}

	for _, user := range mainDB.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return models.User{}, fmt.Errorf("user not found")
}

func loadMainDB() (struct {
	Users     []models.User     `json:"users"`
	Employees []models.Employee `json:"employees"`
}, error) {
	var db struct {
		Users     []models.User     `json:"users"`
		Employees []models.Employee `json:"employees"`
	}

	data, err := os.ReadFile(mainDBFile)
	if err != nil {
		return db, err
	}

	err = json.Unmarshal(data, &db)
	return db, err
}

// ---------- Show Mentor Training Dashboard ----------
func ShowMentorTraining(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to load training data"})
		return
	}

	mainDB, err := loadMainDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to load main data"})
		return
	}

	// Get mentor's batches
	var mentorBatches []models.Batch
	for _, batch := range trainingDB.Batches {
		if batch.MentorID == user.ID {
			mentorBatches = append(mentorBatches, batch)
		}
	}

	// Get students assigned to mentor
	var mentorStudents []map[string]interface{}
	var totalStudents int
	var studentPerformance []map[string]interface{}

	for _, student := range trainingDB.Students {
		if student.MentorID == user.ID {
			totalStudents++

			// Get batch details
			batchName := ""
			duration := ""
			for _, batch := range mentorBatches {
				if batch.ID == student.BatchID {
					batchName = batch.Name
					duration = batch.Duration
					break
				}
			}

			// Calculate attendance
			attendance := calculateStudentAttendance(student.ID, trainingDB.Sessions)

			// Calculate assignment stats
			completedAssignments, totalAssignments := getStudentAssignmentStats(student.ID, trainingDB)

			// Get average marks
			avgMarks := getStudentAverageMarks(student.ID, trainingDB)

			mentorStudents = append(mentorStudents, map[string]interface{}{
				"ID":        student.ID,
				"Name":      student.Name,
				"Course":    student.Course,
				"BatchName": batchName,
				"Duration":  duration,
				"Status":    student.Status,
			})

			studentPerformance = append(studentPerformance, map[string]interface{}{
				"ID":                   student.ID,
				"Name":                 student.Name,
				"Attendance":           attendance,
				"AvgMarks":             avgMarks,
				"CompletedAssignments": completedAssignments,
				"TotalAssignments":     totalAssignments,
				"Remarks":              student.Remarks,
				"Completion":           student.CompletionStatus,
			})
		}
	}

	// Get sessions for mentor's batches
	var mentorSessions []map[string]interface{}
	var totalSessions, completedSessions int

	for _, session := range trainingDB.Sessions {
		for _, batch := range mentorBatches {
			if session.BatchID == batch.ID {
				totalSessions++
				if session.Status == "completed" {
					completedSessions++
				}

				mentorSessions = append(mentorSessions, map[string]interface{}{
					"ID":        session.ID,
					"Date":      session.Date,
					"Topic":     session.Topic,
					"BatchName": batch.Name,
					"Status":    session.Status,
				})
				break
			}
		}
	}

	// Get assignments
	var mentorAssignments []map[string]interface{}
	var pendingAssignments int

	for _, assignment := range trainingDB.Assignments {
		if assignment.CreatedBy == user.ID {
			// Count submissions
			submittedCount := 0
			for _, sub := range trainingDB.Submissions {
				if sub.AssignmentID == assignment.ID {
					submittedCount++
				}
			}

			// Determine target
			target := "Batch"
			totalStudents := 0
			if assignment.StudentID != 0 {
				target = "Individual"
				totalStudents = 1
			} else {
				// Count students in batch
				for _, batch := range mentorBatches {
					if batch.ID == assignment.BatchID {
						totalStudents = len(batch.StudentIDs)
						break
					}
				}
			}

			if submittedCount < totalStudents {
				pendingAssignments++
			}

			mentorAssignments = append(mentorAssignments, map[string]interface{}{
				"ID":             assignment.ID,
				"Title":          assignment.Title,
				"Description":    assignment.Description,
				"Target":         target,
				"DueDate":        assignment.DueDate,
				"SubmittedCount": submittedCount,
				"TotalStudents":  totalStudents,
			})
		}
	}

	// Get mentor name from employees
	mentorName := ""
	for _, emp := range mainDB.Employees {
		if emp.ID == user.ID {
			mentorName = emp.Name
			break
		}
	}

	// Get query parameters for messages
	successMsg := c.Query("success")
	errorMsg := c.Query("error")

	c.HTML(http.StatusOK, "training_mentor.html", gin.H{
		"title":              "Mentor Training Dashboard",
		"Mentor":             gin.H{"Name": mentorName},
		"TotalStudents":      totalStudents,
		"TotalSessions":      totalSessions,
		"CompletedSessions":  completedSessions,
		"PendingAssignments": pendingAssignments,
		"Students":           mentorStudents,
		"Sessions":           mentorSessions,
		"Assignments":        mentorAssignments,
		"StudentPerformance": studentPerformance,
		"Batches":            mentorBatches,
		"Today":              time.Now().Format("2006-01-02"),
		"SuccessMessage":     successMsg,
		"ErrorMessage":       errorMsg,
	})
}

// ---------- Add Student ----------
func AddStudent(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load training data"})
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	course := c.PostForm("course")
	batchID, _ := strconv.Atoi(c.PostForm("batch_id"))

	newStudent := models.Student{
		ID:                   len(trainingDB.Students) + 1,
		Name:                 name,
		Email:                email,
		Phone:                phone,
		Course:               course,
		BatchID:              batchID,
		JoinDate:             time.Now().Format("2006-01-02"),
		Status:               "Active",
		MentorID:             user.ID,
		AttendancePercentage: 0,
		Remarks:              "",
		CompletionStatus:     "ongoing",
	}

	trainingDB.Students = append(trainingDB.Students, newStudent)

	// Add student to batch
	for i, batch := range trainingDB.Batches {
		if batch.ID == batchID {
			trainingDB.Batches[i].StudentIDs = append(trainingDB.Batches[i].StudentIDs, newStudent.ID)
			break
		}
	}

	saveTrainingDB(trainingDB)
	c.Redirect(http.StatusFound, "/training/mentor")
}

// ---------- Add Session ----------
func AddSession(c *gin.Context) {
	_, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load training data"})
		return
	}

	topic := c.PostForm("topic")
	batchID, _ := strconv.Atoi(c.PostForm("batch_id"))
	date := c.PostForm("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	newSession := models.Session{
		ID:         len(trainingDB.Sessions) + 1,
		BatchID:    batchID,
		Date:       date,
		Topic:      topic,
		Trainer:    "",
		Notes:      "",
		Status:     "pending",
		Attendance: make(map[int]bool),
	}

	// Initialize attendance for all students in the batch
	for _, batch := range trainingDB.Batches {
		if batch.ID == batchID {
			for _, studentID := range batch.StudentIDs {
				newSession.Attendance[studentID] = false
			}
			break
		}
	}

	trainingDB.Sessions = append(trainingDB.Sessions, newSession)
	saveTrainingDB(trainingDB)

	c.Redirect(http.StatusFound, "/training/mentor")
}

// ---------- Create Assignment ----------
func CreateAssignment(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load training data"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	assignTo := c.PostForm("assign_to")
	dueDate := c.PostForm("due_date")
	maxScore, _ := strconv.Atoi(c.PostForm("max_score"))

	// Handle file uploads
	var files []string
	form, _ := c.MultipartForm()
	uploadedFiles := form.File["files"]

	// Create uploads directory if not exists
	os.MkdirAll("uploads/assignments", 0755)

	for _, file := range uploadedFiles {
		// Save file to uploads directory
		filename := fmt.Sprintf("uploads/assignments/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, filename); err == nil {
			files = append(files, filename)
		}
	}

	newAssignment := models.Assignment{
		ID:          len(trainingDB.Assignments) + 1,
		Title:       title,
		Description: description,
		CreatedBy:   user.ID,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		DueDate:     dueDate,
		Files:       files,
		MaxScore:    maxScore,
	}

	if assignTo == "batch" {
		batchID, _ := strconv.Atoi(c.PostForm("batch_id"))
		newAssignment.BatchID = batchID
		newAssignment.StudentID = 0
	} else {
		studentID, _ := strconv.Atoi(c.PostForm("student_id"))
		newAssignment.BatchID = 0
		newAssignment.StudentID = studentID
	}

	trainingDB.Assignments = append(trainingDB.Assignments, newAssignment)
	saveTrainingDB(trainingDB)

	c.Redirect(http.StatusFound, "/training/mentor")
}

// ---------- Mark Attendance ----------
func MarkAttendance(c *gin.Context) {
	_, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load training data"})
		return
	}

	sessionID, _ := strconv.Atoi(c.PostForm("session_id"))
	notes := c.PostForm("notes")
	attendance := c.PostFormArray("attendance[]")

	// Update session
	for i, session := range trainingDB.Sessions {
		if session.ID == sessionID {
			// Mark all as false first
			for studentID := range session.Attendance {
				trainingDB.Sessions[i].Attendance[studentID] = false
			}

			// Mark present students
			for _, studentIDStr := range attendance {
				studentID, _ := strconv.Atoi(studentIDStr)
				trainingDB.Sessions[i].Attendance[studentID] = true
			}

			trainingDB.Sessions[i].Notes = notes
			trainingDB.Sessions[i].Status = "completed"
			break
		}
	}

	// Update student attendance percentages
	updateStudentAttendance(&trainingDB)

	saveTrainingDB(trainingDB)
	c.Redirect(http.StatusFound, "/training/mentor")
}

// ---------- Add Batch (supports both form POST and AJAX) ----------
func AddBatch(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		// Check if it's an AJAX request
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.Redirect(http.StatusFound, "/")
		}
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load training data"})
		} else {
			c.Redirect(http.StatusFound, "/training/mentor?error=Failed to load data")
		}
		return
	}

	// Get form values
	name := c.PostForm("name")
	course := c.PostForm("course")
	duration := c.PostForm("duration")
	startDate := c.PostForm("start_date")
	endDate := c.PostForm("end_date")

	// Validate required fields
	if name == "" || course == "" || duration == "" || startDate == "" || endDate == "" {
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		} else {
			c.Redirect(http.StatusFound, "/training/mentor?error=All fields are required")
		}
		return
	}

	// Generate new ID
	newID := 1
	if len(trainingDB.Batches) > 0 {
		newID = trainingDB.Batches[len(trainingDB.Batches)-1].ID + 1
	}

	newBatch := models.Batch{
		ID:         newID,
		Name:       name,
		Course:     course,
		StartDate:  startDate,
		EndDate:    endDate,
		Duration:   duration,
		MentorID:   user.ID,
		StudentIDs: []int{},
		Status:     "active",
	}

	trainingDB.Batches = append(trainingDB.Batches, newBatch)
	err = saveTrainingDB(trainingDB)
	if err != nil {
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save batch"})
		} else {
			c.Redirect(http.StatusFound, "/training/mentor?error=Failed to save batch")
		}
		return
	}

	// Check if it's an AJAX request
	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Batch created successfully",
			"batch": gin.H{
				"ID":        newBatch.ID,
				"Name":      newBatch.Name,
				"Course":    newBatch.Course,
				"Duration":  newBatch.Duration,
				"StartDate": newBatch.StartDate,
				"EndDate":   newBatch.EndDate,
				"Status":    newBatch.Status,
			},
		})
	} else {
		// Regular form submission - redirect back
		c.Redirect(http.StatusFound, "/training/mentor?success=Batch created successfully")
	}
}

// ---------- Initialize Sample Batches ----------
func InitSampleBatches(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load training data"})
		return
	}

	// Check if batches already exist
	if len(trainingDB.Batches) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Batches already exist", "count": len(trainingDB.Batches)})
		return
	}

	// Create sample batches
	sampleBatches := []models.Batch{
		{
			ID:         1,
			Name:       "Python Full Stack Batch 2026-01",
			Course:     "Python Full Stack",
			StartDate:  "2026-01-15",
			EndDate:    "2026-04-15",
			Duration:   "3 months",
			MentorID:   user.ID,
			StudentIDs: []int{},
			Status:     "active",
		},
		{
			ID:         2,
			Name:       "Java Full Stack Batch 2026-01",
			Course:     "Java Full Stack",
			StartDate:  "2026-02-01",
			EndDate:    "2026-05-01",
			Duration:   "3 months",
			MentorID:   user.ID,
			StudentIDs: []int{},
			Status:     "active",
		},
		{
			ID:         3,
			Name:       "Data Science Batch 2026-01",
			Course:     "Data Science",
			StartDate:  "2026-01-20",
			EndDate:    "2026-04-20",
			Duration:   "3 months",
			MentorID:   user.ID,
			StudentIDs: []int{},
			Status:     "active",
		},
	}

	trainingDB.Batches = sampleBatches
	err = saveTrainingDB(trainingDB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save batches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sample batches created successfully",
		"count":   len(sampleBatches),
		"batches": sampleBatches,
	})
}

// ---------- Get All Batches (API) ----------
func GetAllBatches(c *gin.Context) {
	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"batches": trainingDB.Batches,
		"count":   len(trainingDB.Batches),
	})
}

// ---------- Get Session Students (API) ----------
func GetSessionStudents(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	mainDB, err := loadMainDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load main data"})
		return
	}

	var session models.Session
	var batch models.Batch

	for _, s := range trainingDB.Sessions {
		if s.ID == sessionID {
			session = s
			break
		}
	}

	for _, b := range trainingDB.Batches {
		if b.ID == session.BatchID {
			batch = b
			break
		}
	}

	var students []map[string]interface{}
	for _, studentID := range batch.StudentIDs {
		for _, student := range trainingDB.Students {
			if student.ID == studentID {
				// Get student name from employees
				studentName := student.Name
				for _, emp := range mainDB.Employees {
					if emp.ID == studentID {
						studentName = emp.Name
						break
					}
				}

				students = append(students, map[string]interface{}{
					"id":      studentID,
					"name":    studentName,
					"present": session.Attendance[studentID],
				})
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"topic":    session.Topic,
		"students": students,
	})
}

// ---------- Mark Session Complete ----------
func MarkSessionComplete(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	for i, session := range trainingDB.Sessions {
		if session.ID == sessionID {
			trainingDB.Sessions[i].Status = "completed"
			break
		}
	}

	saveTrainingDB(trainingDB)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ---------- Add Session Notes ----------
func AddSessionNotes(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Notes string `json:"notes"`
	}
	c.BindJSON(&req)

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	for i, session := range trainingDB.Sessions {
		if session.ID == sessionID {
			trainingDB.Sessions[i].Notes = req.Notes
			break
		}
	}

	saveTrainingDB(trainingDB)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ---------- Get Assignment Submissions (API) ----------
func GetAssignmentSubmissions(c *gin.Context) {
	assignmentID, _ := strconv.Atoi(c.Param("id"))

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	mainDB, err := loadMainDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load main data"})
		return
	}

	var submissions []map[string]interface{}

	// Create grade map
	gradeMap := make(map[int]models.Grade)
	for _, grade := range trainingDB.Grades {
		gradeMap[grade.SubmissionID] = grade
	}

	for _, sub := range trainingDB.Submissions {
		if sub.AssignmentID == assignmentID {
			// Get student name
			studentName := ""
			for _, student := range trainingDB.Students {
				if student.ID == sub.StudentID {
					studentName = student.Name
					break
				}
			}
			// Try to get from employees if not found
			if studentName == "" {
				for _, emp := range mainDB.Employees {
					if emp.ID == sub.StudentID {
						studentName = emp.Name
						break
					}
				}
			}

			grade, exists := gradeMap[sub.ID]
			submission := map[string]interface{}{
				"id":              sub.ID,
				"student_name":    studentName,
				"submitted_at":    sub.SubmittedAt,
				"files":           sub.Files,
				"late_submission": sub.LateSubmission,
				"graded":          exists,
			}

			if exists {
				submission["score"] = grade.Score
				submission["feedback"] = grade.Feedback
			}

			submissions = append(submissions, submission)
		}
	}

	c.JSON(http.StatusOK, gin.H{"submissions": submissions})
}

// ---------- Grade Assignment ----------
func GradeAssignment(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		SubmissionID int    `json:"submission_id"`
		AssignmentID int    `json:"assignment_id"`
		Score        int    `json:"score"`
		Feedback     string `json:"feedback"`
	}
	c.BindJSON(&req)

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	newGrade := models.Grade{
		ID:           len(trainingDB.Grades) + 1,
		SubmissionID: req.SubmissionID,
		AssignmentID: req.AssignmentID,
		StudentID:    0, // Will be filled from submission
		Score:        req.Score,
		Feedback:     req.Feedback,
		GradedAt:     time.Now().Format("2006-01-02 15:04:05"),
		GradedBy:     user.ID,
	}

	// Get student ID from submission
	for _, sub := range trainingDB.Submissions {
		if sub.ID == req.SubmissionID {
			newGrade.StudentID = sub.StudentID
			break
		}
	}

	trainingDB.Grades = append(trainingDB.Grades, newGrade)

	// Update submission status
	for i, sub := range trainingDB.Submissions {
		if sub.ID == req.SubmissionID {
			trainingDB.Submissions[i].Status = "graded"
			break
		}
	}

	saveTrainingDB(trainingDB)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ---------- Add Student Remarks ----------
func AddStudentRemarks(c *gin.Context) {
	studentID, _ := strconv.Atoi(c.PostForm("student_id"))
	remarks := c.PostForm("remarks")
	completion := c.PostForm("completion")

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	for i, student := range trainingDB.Students {
		if student.ID == studentID {
			trainingDB.Students[i].Remarks = remarks
			trainingDB.Students[i].CompletionStatus = completion
			break
		}
	}

	saveTrainingDB(trainingDB)
	c.Redirect(http.StatusFound, "/training/mentor")
}

// ---------- Get Student Details (API) ----------
func GetStudentDetails(c *gin.Context) {
	studentID, _ := strconv.Atoi(c.Param("id"))

	trainingDB, err := loadTrainingDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load data"})
		return
	}

	mainDB, err := loadMainDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load main data"})
		return
	}

	for _, student := range trainingDB.Students {
		if student.ID == studentID {
			// Get name from employees if available
			name := student.Name
			for _, emp := range mainDB.Employees {
				if emp.ID == studentID {
					name = emp.Name
					break
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"id":   student.ID,
				"name": name,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
}

// ---------- Helper Functions ----------
func updateStudentAttendance(db *models.TrainingDB) {
	// Calculate attendance for each student
	for i, student := range db.Students {
		totalSessions := 0
		presentSessions := 0

		for _, session := range db.Sessions {
			if session.BatchID == student.BatchID {
				totalSessions++
				if session.Attendance[student.ID] {
					presentSessions++
				}
			}
		}

		if totalSessions > 0 {
			db.Students[i].AttendancePercentage = float64(presentSessions) / float64(totalSessions) * 100
		}
	}
}

func calculateStudentAttendance(studentID int, sessions []models.Session) float64 {
	totalSessions := 0
	presentSessions := 0

	for _, session := range sessions {
		if present, exists := session.Attendance[studentID]; exists {
			totalSessions++
			if present {
				presentSessions++
			}
		}
	}

	if totalSessions == 0 {
		return 0
	}
	return float64(presentSessions) / float64(totalSessions) * 100
}

func getStudentAssignmentStats(studentID int, db models.TrainingDB) (completed int, total int) {
	// Get assignments for this student
	assignmentMap := make(map[int]bool)

	for _, assignment := range db.Assignments {
		if assignment.StudentID == studentID || (assignment.StudentID == 0 && batchContainsStudent(assignment.BatchID, studentID, db)) {
			assignmentMap[assignment.ID] = true
		}
	}

	total = len(assignmentMap)

	// Count submissions
	for _, sub := range db.Submissions {
		if sub.StudentID == studentID && assignmentMap[sub.AssignmentID] {
			completed++
		}
	}

	return
}

func getStudentAverageMarks(studentID int, db models.TrainingDB) float64 {
	var totalMarks int
	var count int

	for _, grade := range db.Grades {
		if grade.StudentID == studentID {
			totalMarks += grade.Score
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return float64(totalMarks) / float64(count)
}

func batchContainsStudent(batchID int, studentID int, db models.TrainingDB) bool {
	for _, batch := range db.Batches {
		if batch.ID == batchID {
			for _, id := range batch.StudentIDs {
				if id == studentID {
					return true
				}
			}
		}
	}
	return false
}

// ---------- Redirect Handlers (for backward compatibility) ----------
func RedirectMentorTraining(c *gin.Context) {
	c.Redirect(http.StatusFound, "/training/mentor")
}

func RedirectStudentTraining(c *gin.Context) {
	c.Redirect(http.StatusFound, "/training/student")
}
