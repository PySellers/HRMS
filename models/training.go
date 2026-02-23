package models

type Student struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Course    string `json:"course"`
	BatchID   int    `json:"batch_id"`
	JoinDate  string `json:"join_date"`
	Status    string `json:"status"` // Active / Completed / Dropped
}

type Batch struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Mentor       string `json:"mentor"`
	CourseName   string `json:"course_name"`
	Description  string `json:"description"`
}

type Session struct {
	ID        int    `json:"id"`
	BatchID   int    `json:"batch_id"`
	Date      string `json:"date"`
	Topic     string `json:"topic"`
	Trainer   string `json:"trainer"`
	Attendance map[int]bool `json:"attendance"` // studentID → true/false
}

type Material struct {
	ID        int    `json:"id"`
	BatchID   int    `json:"batch_id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Type      string `json:"type"` // Video / PDF / Link / Assignment
	Uploaded  string `json:"uploaded"`
}

type Evaluation struct {
	ID        int    `json:"id"`
	StudentID int    `json:"student_id"`
	Date      string `json:"date"`
	Topic     string `json:"topic"`
	Grade     string `json:"grade"`
	Feedback  string `json:"feedback"`
}

type TrainingDB struct {
	Students    []Student    `json:"students"`
	Batches     []Batch      `json:"batches"`
	Sessions    []Session    `json:"sessions"`
	Materials   []Material   `json:"materials"`
	Evaluations []Evaluation `json:"evaluations"`
}
