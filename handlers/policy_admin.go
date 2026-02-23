package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

// Show all company policies (admin view)
func ShowPolicies(c *gin.Context) {
	data, _ := os.ReadFile(dbFile)
	var db models.DB
	json.Unmarshal(data, &db)

	c.HTML(http.StatusOK, "policies.html", gin.H{
		"policies": db.Policies,
	})
}

// Add new policy
func AddPolicy(c *gin.Context) {
	title := c.PostForm("title")

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "File upload failed: %v", err)
		return
	}

	filename := filepath.Base(file.Filename)
	savePath := filepath.Join("static", "docs", filename)

	// ensure folder exists
	os.MkdirAll("static/docs", os.ModePerm)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.String(http.StatusInternalServerError, "Save failed: %v", err)
		return
	}

	data, _ := os.ReadFile(dbFile)
	var db models.DB
	json.Unmarshal(data, &db)

	newPolicy := models.Policy{
		ID:    len(db.Policies) + 1,
		Title: title,
		URL:   "/static/docs/" + filename,
	}

	db.Policies = append(db.Policies, newPolicy)

	saveDB(db)
	c.Redirect(http.StatusFound, "/admin/policies")
}
