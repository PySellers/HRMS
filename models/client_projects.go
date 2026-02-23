package models

type Client struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Company string `json:"company"`
	Address string `json:"address"`
	Notes   string `json:"notes"`
}

type ClientProject struct {
	ID          int     `json:"id"`
	ClientID    int     `json:"client_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"` // Planned / In Progress / Completed
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Budget      float64 `json:"budget"`
	HoursLogged float64 `json:"hours_logged"`
}

type Task struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"` // To Do / In Progress / Review / Done
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

type Invoice struct {
	ID        int     `json:"id"`
	ClientID  int     `json:"client_id"`
	ProjectID int     `json:"project_id"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
	Status    string  `json:"status"` // Paid / Pending / Overdue
}

type ClientDB struct {
	Clients  []Client        `json:"clients"`
	Projects []ClientProject `json:"projects"`
	Tasks    []Task          `json:"tasks"`
	Invoices []Invoice       `json:"invoices"`
}
