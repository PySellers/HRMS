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

// Show finance dashboard
func ShowFinanceDashboard(c *gin.Context) {
	data, _ := os.ReadFile(dbFile)
	var db models.DB
	json.Unmarshal(data, &db)

	totalIncome := 0.0
	totalExpense := 0.0
	for _, f := range db.Finance {
		if f.Type == "income" {
			totalIncome += f.Amount
		} else if f.Type == "expense" {
			totalExpense += f.Amount
		}
	}

	balance := totalIncome - totalExpense

	c.HTML(http.StatusOK, "finance.html", gin.H{
		"finance":      db.Finance,
		"totalIncome":  totalIncome,
		"totalExpense": totalExpense,
		"balance":      balance,
	})
}

// Add a new finance entry
func AddFinance(c *gin.Context) {
	data, _ := os.ReadFile(dbFile)
	var db models.DB
	json.Unmarshal(data, &db)

	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)

	newFinance := models.Finance{
		ID:          len(db.Finance) + 1,
		Type:        c.PostForm("type"),
		Category:    c.PostForm("category"),
		Description: c.PostForm("description"),
		Amount:      amount,
		Date:        time.Now().Format("2006-01-02"),
	}

	db.Finance = append(db.Finance, newFinance)
	saveDB(db)

	c.Redirect(http.StatusFound, "/admin/finance")
}
