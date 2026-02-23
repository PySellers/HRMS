package models

type Incident struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`   // Hardware, Network, Software
	Priority    string `json:"priority"`   // Low, Medium, High, Critical
	Status      string `json:"status"`     // Open, In Progress, Resolved, Closed
	AssignedTo  string `json:"assigned_to"`
	RootCause   string `json:"root_cause"`
	Resolution  string `json:"resolution"`
	CreatedDate string `json:"created_date"`
	ClosedDate  string `json:"closed_date"`
}

type Asset struct {
	ID           int    `json:"id"`
	Type         string `json:"type"`        // Laptop, Server, License, Router
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	Owner        string `json:"owner"`
	Status       string `json:"status"`      // Active, In Maintenance, Retired
	PurchaseDate string `json:"purchase_date"`
	WarrantyTill string `json:"warranty_till"`
}

type ServiceCatalog struct {
	ID          int    `json:"id"`
	ServiceName string `json:"service_name"`
	Description string `json:"description"`
}

type ServiceRequest struct {
	ID          int    `json:"id"`
	RequestedBy string `json:"requested_by"`
	ServiceID   int    `json:"service_id"`
	Status      string `json:"status"`
	RequestDate string `json:"request_date"`
}

type ITSMDB struct {
	Incidents      []Incident      `json:"incidents"`
	Assets         []Asset         `json:"assets"`
	Services       []ServiceCatalog `json:"services"`
	ServiceRequests []ServiceRequest `json:"service_requests"`
}
