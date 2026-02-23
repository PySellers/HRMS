package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"pysellers-erp-go/models"
	"pysellers-erp-go/utils"

	"github.com/gin-gonic/gin"
)

// ================================
// LOAD EMPLOYEES
// ================================
func loadEmployees() ([]models.Employee, error) {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		return nil, err
	}

	var db map[string]interface{}
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}

	employeesData, _ := json.Marshal(db["employees"])
	var employees []models.Employee
	json.Unmarshal(employeesData, &employees)

	return employees, nil
}

// ================================
// SAVE EMPLOYEES
// ================================
func saveEmployees(employees []models.Employee) error {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		return err
	}

	var db map[string]interface{}
	json.Unmarshal(data, &db)
	db["employees"] = employees

	newData, _ := json.MarshalIndent(db, "", "  ")
	return os.WriteFile(dbFile, newData, 0644)
}

// ================================
// SHOW EMPLOYEE LIST (HTML)
// ================================
func ShowEmployees(c *gin.Context) {
	employees, err := loadEmployees()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load employees"})
		return
	}

	c.HTML(http.StatusOK, "employees.html", gin.H{
		"employees": employees,
	})
}

// ================================
// ADD EMPLOYEE (FORM HANDLER)
// ================================
func AdminAddEmployee(c *gin.Context) {
	// Read full DB to perform atomic checks and writes
	db, err := utils.ReadDB()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load database"})
		return
	}

	// ---------- FORM DATA ----------
	empID := strings.TrimSpace(strings.ToUpper(c.PostForm("employee_id")))
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	phone := strings.TrimSpace(c.PostForm("phone"))
	address := strings.TrimSpace(c.PostForm("address"))
	projectName := strings.TrimSpace(c.PostForm("project_name"))
	projectManager := strings.TrimSpace(c.PostForm("project_manager"))

	// coworkers is ONE input (comma separated)
	coworkers := strings.Split(c.PostForm("coworkers"), ",")
	for i := range coworkers {
		coworkers[i] = strings.TrimSpace(coworkers[i])
	}

	// Check duplicates: employeeID, email or existing username
	if utils.EmployeeOrUserExists(db, empID, email) {
		// render page with existing employees for admin view
		employees, _ := loadEmployees()
		c.HTML(http.StatusConflict, "employees.html", gin.H{
			"employees": employees,
			"error":     "Employee ID or email already exists. No new account created.",
		})
		return
	}

	// Generate password
	rawPassword := utils.GeneratePassword()
	hashedPassword, _ := utils.HashPassword(rawPassword)

	// ---------- CREATE EMPLOYEE + USER in central DB
	newEmp := models.Employee{
		ID:         utils.NextEmployeeID(db),
		EmployeeID: empID,
		Name:       name,
		Email:      email,
		Phone:      phone,
		Address:    address,
		Project: models.Project{
			Name:      projectName,
			Manager:   projectManager,
			Coworkers: coworkers,
		},
		Skills:         []string{},
		Certifications: []string{},
		ProfilePic:     "",
		IsActive:       true,
	}

	db.Employees = append(db.Employees, newEmp)

	// Use employee ID as username for admin-created accounts
	username := empID
	newUser := models.User{
		ID:       utils.NextUserID(db),
		Username: username,
		Password: hashedPassword,
		Role:     "employee",
	}
	db.Users = append(db.Users, newUser)

	// Save DB once
	if err := utils.WriteDB(db); err != nil {
		c.HTML(http.StatusInternalServerError, "employees.html", gin.H{
			"error": "Failed to save data",
		})
		return
	}

	// Send credentials (log in background)
	log.Println("📧 Sending credentials to:", email)
	go func() {
		if err := utils.SendCredentials(email, name, username, rawPassword, "employee"); err != nil {
			log.Println("❌ Email send FAILED:", err)
		} else {
			log.Println("✅ Email sent successfully to:", email)
		}
	}()

	c.Redirect(302, "/admin/employees")
}

// ================================
// UPDATE USERS SECTION
// ================================
func updateUsers(username, password, role string) {
	data, _ := os.ReadFile(dbFile)

	var db map[string]interface{}
	json.Unmarshal(data, &db)

	usersData, _ := json.Marshal(db["users"])
	var users []models.User
	json.Unmarshal(usersData, &users)

	newUser := models.User{
		ID:       len(users) + 1,
		Username: username,
		Password: password, // already hashed
		Role:     role,
	}

	users = append(users, newUser)
	db["users"] = users

	newData, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(dbFile, newData, 0644)
}

// ================================
// DELETE EMPLOYEE
// ================================
func DeleteEmployee(c *gin.Context) {
	employees, _ := loadEmployees()
	id, _ := strconv.Atoi(c.Param("id"))

	newList := []models.Employee{}
	var deletedUsername string

	for _, e := range employees {
		if e.ID == id {
			deletedUsername = e.EmployeeID
			continue
		}
		newList = append(newList, e)
	}

	saveEmployees(newList)

	if deletedUsername != "" {
		cleanupUser(deletedUsername)
	}

	c.Redirect(302, "/admin/employees")
}

// ================================
// EDIT EMPLOYEE (HTML)
// ================================
func EditEmployee(c *gin.Context) {
	employees, _ := loadEmployees()
	id, _ := strconv.Atoi(c.Param("id"))

	var emp models.Employee
	for _, e := range employees {
		if e.ID == id {
			emp = e
			break
		}
	}

	c.HTML(http.StatusOK, "employee_edit.html", gin.H{
		"employee": emp,
	})
}

// ================================
// UPDATE EMPLOYEE
// ================================
func UpdateEmployee(c *gin.Context) {
	employees, _ := loadEmployees()
	id, _ := strconv.Atoi(c.Param("id"))

	for i, e := range employees {
		if e.ID == id {
			employees[i].EmployeeID = c.PostForm("employee_id")
			employees[i].Name = c.PostForm("name")
			employees[i].Email = c.PostForm("email")
			employees[i].Phone = c.PostForm("phone")
			employees[i].Address = c.PostForm("address")
			employees[i].Project.Name = c.PostForm("project_name")
			employees[i].Project.Manager = c.PostForm("project_manager")

			coworkers := strings.Split(c.PostForm("coworkers"), ",")
			for i := range coworkers {
				coworkers[i] = strings.TrimSpace(coworkers[i])
			}
			employees[i].Project.Coworkers = coworkers
			break
		}
	}

	saveEmployees(employees)
	c.Redirect(302, "/admin/employees")
}

// ================================
// REMOVE USER
// ================================
func cleanupUser(username string) {
	data, _ := os.ReadFile(dbFile)

	var db map[string]interface{}
	json.Unmarshal(data, &db)

	usersData, _ := json.Marshal(db["users"])
	var users []models.User
	json.Unmarshal(usersData, &users)

	filtered := []models.User{}
	for _, u := range users {
		if u.Username != username {
			filtered = append(filtered, u)
		}
	}

	db["users"] = filtered
	newData, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(dbFile, newData, 0644)
}
