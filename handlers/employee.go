package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"pysellers-erp-go/models"
	"pysellers-erp-go/utils"
)

// ================================
// TOGGLE EMPLOYEE ACTIVE / INACTIVE
// ================================
func ToggleEmployeeStatus(c *gin.Context) {

	utils.DBMutex.Lock()
	defer utils.DBMutex.Unlock()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/employees")
		return
	}

	db, err := utils.ReadDB()
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/employees")
		return
	}

	for i := range db.Employees {
		if db.Employees[i].ID == id {
			db.Employees[i].IsActive = !db.Employees[i].IsActive
			break
		}
	}

	_ = utils.WriteDB(db)
	c.Redirect(http.StatusFound, "/admin/employees")
}

// ================================
// ADD NEW EMPLOYEE (ADMIN)
// ================================
func AddEmployee(c *gin.Context) {

	utils.DBMutex.Lock()
	defer utils.DBMutex.Unlock()

	// -------------------------------
	// READ FORM VALUES
	// -------------------------------
	employeeID := strings.ToUpper(strings.TrimSpace(c.PostForm("employee_id")))
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	phone := c.PostForm("phone")
	address := c.PostForm("address")
	role := strings.TrimSpace(c.PostForm("role"))

	projectName := c.PostForm("project_name")
	projectManager := c.PostForm("project_manager")
	coworkersRaw := c.PostForm("coworkers")

	// -------------------------------
	// BASIC VALIDATION
	// -------------------------------
	if employeeID == "" || email == "" || role == "" {
		db, _ := utils.ReadDB()
		c.HTML(http.StatusBadRequest, "employees.html", gin.H{
			"employees": db.Employees,
			"error":     "Employee ID, Email and Role are required.",
		})
		return
	}

	// -------------------------------
	// READ DATABASE
	// -------------------------------
	db, err := utils.ReadDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "employees.html", gin.H{
			"error": "Database error. Please try again.",
		})
		return
	}

	// -------------------------------
	// DUPLICATE CHECK
	// -------------------------------
	if utils.EmployeeOrUserExists(db, employeeID, email) {
		c.HTML(http.StatusConflict, "employees.html", gin.H{
			"employees": db.Employees,
			"error":     "Employee ID or Email already exists.",
		})
		return
	}

	// -------------------------------
	// CREATE EMPLOYEE
	// -------------------------------
	newEmp := models.Employee{
		ID:         utils.NextEmployeeID(db),
		EmployeeID: employeeID,
		Name:       name,
		Email:      email,
		Phone:      phone,
		Role:       role,
		Address:    address,
		IsActive:   true,
	}

	// Project only if NOT client
	if role != "client" {
		newEmp.Project = models.Project{
			Name:      projectName,
			Manager:   projectManager,
			Coworkers: strings.Split(coworkersRaw, ","),
		}
	}

	db.Employees = append(db.Employees, newEmp)

	// -------------------------------
	// CREATE USER ACCOUNT
	// -------------------------------
	rawPassword := utils.GeneratePassword()
	hashedPassword, err := utils.HashPassword(rawPassword)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "employees.html", gin.H{
			"employees": db.Employees,
			"error":     "Password generation failed.",
		})
		return
	}

	user := models.User{
		ID:       utils.NextUserID(db),
		Username: employeeID,
		Password: hashedPassword,
		Role:     role,
	}

	db.Users = append(db.Users, user)

	// -------------------------------
	// SAVE DATABASE
	// -------------------------------
	if err := utils.WriteDB(db); err != nil {
		c.HTML(http.StatusInternalServerError, "employees.html", gin.H{
			"employees": db.Employees,
			"error":     "Failed to save employee data.",
		})
		return
	}

	// -------------------------------
	// SEND EMAIL
	// -------------------------------
	log.Println("📧 Sending credentials to:", email)

	go utils.SendCredentials(
		email,
		name,
		employeeID,
		rawPassword,
		role,
	)

	// -------------------------------
	// REDIRECT
	// -------------------------------
	c.Redirect(http.StatusFound, "/admin/employees")
}
