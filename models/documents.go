package models

type Document struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	FileName    string `json:"file_name"`
	FilePath    string `json:"file_path"`
	Category    string `json:"category"` // Certificate / Financial / Agreement / Audit
	Version     string `json:"version"`
	UploadedBy  string `json:"uploaded_by"`
	UploadedAt  string `json:"uploaded_at"`
	LastViewed  string `json:"last_viewed"`
	ViewCount   int    `json:"view_count"`
	Description string `json:"description"`
}

type DocumentDB struct {
	Documents []Document `json:"documents"`
}
