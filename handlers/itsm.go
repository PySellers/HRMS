package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pysellers-erp-go/models"
)

var itsmDBFile = "data/itsm.json"

// Dashboard
func ShowITSM(c *gin.Context) {
	data, _ := os.ReadFile(itsmDBFile)
	var db models.ITSMDB
	json.Unmarshal(data, &db)
	c.HTML(http.StatusOK, "itsm.html", gin.H{
		"incidents": db.Incidents,
		"assets": db.Assets,
		"services": db.Services,
		"requests": db.ServiceRequests,
	})
}

// Add Incident
func AddIncident(c *gin.Context) {
	data, _ := os.ReadFile(itsmDBFile)
	var db models.ITSMDB
	json.Unmarshal(data, &db)

	newIncident := models.Incident{
		ID: len(db.Incidents) + 1,
		Title: c.PostForm("title"),
		Description: c.PostForm("description"),
		Category: c.PostForm("category"),
		Priority: c.PostForm("priority"),
		Status: "Open",
		CreatedDate: time.Now().Format("2006-01-02"),
	}
	db.Incidents = append(db.Incidents, newIncident)
	saveITSM(db)
	c.Redirect(http.StatusFound, "/admin/itsm")
}

// Add Asset
func AddAsset(c *gin.Context) {
	data, _ := os.ReadFile(itsmDBFile)
	var db models.ITSMDB
	json.Unmarshal(data, &db)

	newAsset := models.Asset{
		ID: len(db.Assets) + 1,
		Type: c.PostForm("type"),
		Name: c.PostForm("name"),
		SerialNumber: c.PostForm("serial_number"),
		Owner: c.PostForm("owner"),
		Status: "Active",
		PurchaseDate: c.PostForm("purchase_date"),
		WarrantyTill: c.PostForm("warranty_till"),
	}
	db.Assets = append(db.Assets, newAsset)
	saveITSM(db)
	c.Redirect(http.StatusFound, "/admin/itsm")
}

// Add Service
func AddService(c *gin.Context) {
	data, _ := os.ReadFile(itsmDBFile)
	var db models.ITSMDB
	json.Unmarshal(data, &db)

	newService := models.ServiceCatalog{
		ID: len(db.Services) + 1,
		ServiceName: c.PostForm("service_name"),
		Description: c.PostForm("description"),
	}
	db.Services = append(db.Services, newService)
	saveITSM(db)
	c.Redirect(http.StatusFound, "/admin/itsm")
}

// Add Service Request
func AddServiceRequest(c *gin.Context) {
	data, _ := os.ReadFile(itsmDBFile)
	var db models.ITSMDB
	json.Unmarshal(data, &db)

	serviceID, _ := strconv.Atoi(c.PostForm("service_id"))
	newRequest := models.ServiceRequest{
		ID: len(db.ServiceRequests) + 1,
		RequestedBy: c.PostForm("requested_by"),
		ServiceID: serviceID,
		Status: "Pending",
		RequestDate: time.Now().Format("2006-01-02"),
	}
	db.ServiceRequests = append(db.ServiceRequests, newRequest)
	saveITSM(db)
	c.Redirect(http.StatusFound, "/admin/itsm")
}

// Save helper
func saveITSM(db models.ITSMDB) {
	out, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(itsmDBFile, out, 0644)
}
