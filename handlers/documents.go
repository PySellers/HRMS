package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

var docsFile = "data/documents.json"

// Show Documents Page
func ShowDocuments(c *gin.Context) {
	data, _ := os.ReadFile(docsFile)
	var db models.DocumentDB
	json.Unmarshal(data, &db)

	c.HTML(http.StatusOK, "documents.html", gin.H{"documents": db.Documents})
}

// Upload Document
func UploadDocument(c *gin.Context) {
	title := c.PostForm("title")
	category := c.PostForm("category")
	version := c.PostForm("version")
	desc := c.PostForm("description")

	file, err := c.FormFile("document")
	if err != nil {
		c.String(http.StatusBadRequest, "Upload failed: %v", err)
		return
	}

	saveDir := "static/docs"
	os.MkdirAll(saveDir, 0755)
	filename := filepath.Base(file.Filename)
	savePath := filepath.Join(saveDir, filename)
	c.SaveUploadedFile(file, savePath)

	data, _ := os.ReadFile(docsFile)
	var db models.DocumentDB
	json.Unmarshal(data, &db)

	newDoc := models.Document{
		ID:          len(db.Documents) + 1,
		Title:       title,
		FileName:    filename,
		FilePath:    savePath,
		Category:    category,
		Version:     version,
		UploadedBy:  "Admin",
		UploadedAt:  time.Now().Format("2006-01-02 15:04"),
		Description: desc,
	}
	db.Documents = append(db.Documents, newDoc)

	out, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(docsFile, out, 0644)
	c.Redirect(http.StatusFound, "/admin/documents")
}
