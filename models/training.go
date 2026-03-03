package models
 
type TrainingDB struct {
    Students    []Student    `json:"students"`
    Batches     []Batch      `json:"batches"`
    Sessions    []Session    `json:"sessions"`
    Materials   []Material   `json:"materials"`
    Assignments []Assignment `json:"assignments"`
    Submissions []Submission `json:"submissions"`
    Grades      []Grade      `json:"grades"`
    Fees        []Fee        `json:"fees,omitempty"`
    Certificates []Certificate `json:"certificates,omitempty"`
    Courses     []Course     `json:"courses,omitempty"`
}
 
type Student struct {
    ID                 int     `json:"id"`
    Name               string  `json:"name"`
    Email              string  `json:"email"`
    Phone              string  `json:"phone"`
    Course             string  `json:"course"`
    BatchID            int     `json:"batch_id"`
    JoinDate           string  `json:"join_date"`
    Status             string  `json:"status"`
    MentorID           int     `json:"mentor_id"`
    AttendancePercentage float64 `json:"attendance_percentage"` // This is the field name
    Remarks            string  `json:"remarks"`
    CompletionStatus   string  `json:"completion_status"` // This is the field name
}
 
type Batch struct {
    ID         int      `json:"id"`
    Name       string   `json:"name"`
    Course     string   `json:"course"`
    StartDate  string   `json:"start_date"`
    EndDate    string   `json:"end_date"`
    Duration   string   `json:"duration"`
    MentorID   int      `json:"mentor_id"`
    StudentIDs []int    `json:"student_ids"`
    Status     string   `json:"status"`
}
 
type Session struct {
    ID         int            `json:"id"`
    BatchID    int            `json:"batch_id"`
    Date       string         `json:"date"`
    Topic      string         `json:"topic"`
    Trainer    string         `json:"trainer"`
    Notes      string         `json:"notes"`
    Status     string         `json:"status"`
    Attendance map[int]bool   `json:"attendance"`
}
 
type Material struct {
    ID       int    `json:"id"`
    BatchID  int    `json:"batch_id"`
    Title    string `json:"title"`
    URL      string `json:"url"`
    Type     string `json:"type"`
    Uploaded string `json:"uploaded"`
}
 
type Assignment struct {
    ID          int      `json:"id"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    BatchID     int      `json:"batch_id"`
    StudentID   int      `json:"student_id,omitempty"`
    CreatedBy   int      `json:"created_by"`
    CreatedAt   string   `json:"created_at"`
    DueDate     string   `json:"due_date"`
    Files       []string `json:"files"`
    MaxScore    int      `json:"max_score"`
}
 
type Submission struct {
    ID           int      `json:"id"`
    AssignmentID int      `json:"assignment_id"`
    StudentID    int      `json:"student_id"`
    SubmittedAt  string   `json:"submitted_at"`
    Files        []string `json:"files"`
    Status       string   `json:"status"`
    LateSubmission bool   `json:"late_submission"`
}
 
type Grade struct {
    ID           int    `json:"id"`
    SubmissionID int    `json:"submission_id"`
    AssignmentID int    `json:"assignment_id"`
    StudentID    int    `json:"student_id"`
    Score        int    `json:"score"`
    Feedback     string `json:"feedback"`
    GradedAt     string `json:"graded_at"`
    GradedBy     int    `json:"graded_by"`
}
 
type Fee struct {
    ID           int     `json:"id"`
    StudentID    int     `json:"student_id"`
    TotalFee     float64 `json:"total_fee"`
    PaidAmount   float64 `json:"paid_amount"`
    PendingAmount float64 `json:"pending_amount"`
    DueDate      string  `json:"due_date"`
    Status       string  `json:"status"`
}
 
type Certificate struct {
    ID            int    `json:"id"`
    StudentID     int    `json:"student_id"`
    Course        string `json:"course"`
    IssueDate     string `json:"issue_date"`
    CertificateURL string `json:"certificate_url"`
    IsEligible    bool   `json:"is_eligible"`
    IsIssued      bool   `json:"is_issued"`
}
 
type Course struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Duration    string `json:"duration"`
    Curriculum  string `json:"curriculum"`
    NextCourse  string `json:"next_course"`
}
 