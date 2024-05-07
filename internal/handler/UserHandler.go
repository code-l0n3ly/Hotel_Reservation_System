package Handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db              *sql.DB
	UserIdReference int64
	cache           map[string]Entities.User // Cache to hold users in memory
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		db:              db,
		UserIdReference: 0,
		cache:           make(map[string]Entities.User),
	}
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (UserHandler *UserHandler) LoadUsersIntoCache() error {
	query := `
        SELECT 
            u.UserID, u.AddressID, u.Name, u.PhoneNumber, u.Email, u.Password, u.CreateTime, u.UserRole,
            a.AddressID, a.Country, a.City, a.State, a.Street, a.PostalCode, a.AdditionalNumber, a.MapLocation, a.Latitude, a.Longitude
        FROM 
            User u
        LEFT JOIN 
            Address a ON u.AddressID = a.AddressID
    `
	rows, err := UserHandler.db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var user Entities.User
		var address Entities.Address
		if err := rows.Scan(&user.UserID, &user.AddressID, &user.Name, &user.PhoneNumber, &user.Email, &user.Password, &createTime, &user.UserRole, &address.AddressID, &address.Country, &address.City, &address.State, &address.Street, &address.PostalCode, &address.AdditionalNumber, &address.MapLocation, &address.Latitude, &address.Longitude); err != nil {
			fmt.Println(err.Error())
			return err
		}
		user.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))

		// Save the grouped object to the cache
		user.Address = address
		fmt.Println(user)
		UserHandler.cache[user.UserID] = user
	}
	return rows.Err()
}

func (UserHandler *UserHandler) GetUserByID(userID string) (Entities.User, bool) {
	user, exists := UserHandler.cache[userID]
	return user, exists
}

func (UserHandler *UserHandler) CreateUserHandler(c *gin.Context) {
	UserHandler.LoadUsersIntoCache()
	var user Entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists in the cache
	if !user.IsEmailValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is not valid"})
		return
	}
	// else if !user.IsPasswordStrong() {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong enough"})
	// 	return
	// }
	for _, existingUser := range UserHandler.cache {
		if existingUser.Email == user.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
	}

	tx, err := UserHandler.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create property"})
		return
	}

	addressQuery := `INSERT INTO Address (Country, City, State, Street, PostalCode, AdditionalNumber, MapLocation, Latitude, Longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	addressResult, err := tx.Exec(addressQuery, user.Address.Country, user.Address.City, user.Address.State, user.Address.Street, user.Address.PostalCode, user.Address.AdditionalNumber, user.Address.MapLocation, user.Address.Latitude, user.Address.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create address" + err.Error()})
		return
	}
	AddressID, _ := addressResult.LastInsertId()

	user.AddressID = strconv.FormatInt(AddressID, 10)

	query := `INSERT INTO User (AddressID, Name, Email, PhoneNumber, Password, UserRole) VALUES (?, ?, ?, ?, ?, ?)`
	r, err := tx.Exec(query, AddressID, user.Name, user.Email, user.PhoneNumber, user.Password, user.UserRole)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
	}
	tx.Commit()
	// Add the new user to the cache
	id, _ := r.LastInsertId()
	user.UserID = strconv.FormatInt(id, 10)
	UserHandler.LoadUsersIntoCache()
	response := Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    UserHandler.cache[user.UserID],
	}

	c.JSON(http.StatusCreated, response)
}

func (UserHandler *UserHandler) GetUserHandler(c *gin.Context) {
	UserHandler.LoadUsersIntoCache()
	userID := c.Param("id")

	// Try to get the user from the cache
	user, ok := UserHandler.cache[userID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := Response{
		Status:  "success",
		Message: "User retrieved successfully",
		Data:    user,
	}
	c.JSON(http.StatusOK, response)
}

func (UserHandler *UserHandler) GetUsersHandler(c *gin.Context) {
	UserHandler.LoadUsersIntoCache()
	var users []Entities.User
	for _, user := range UserHandler.cache {
		users = append(users, user)
	}

	response := Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    users,
	}

	c.JSON(http.StatusOK, response)
}

func (UserHandler *UserHandler) UpdateUserHandler(c *gin.Context) {
	userID := c.Param("id")
	UserHandler.LoadUsersIntoCache()
	oldUser, exists := UserHandler.GetUserByID(userID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var newUser Entities.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newUser.Name == "" {
		newUser.Name = oldUser.Name
	}
	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}
	if newUser.UserRole == "" {
		newUser.UserRole = oldUser.UserRole
	}
	if newUser.PhoneNumber == "" {
		newUser.PhoneNumber = oldUser.PhoneNumber
	}
	query := `UPDATE User SET Name = ?, PhoneNumber = ?, Email = ?, UserRole = ? WHERE UserID = ?`
	_, err := UserHandler.db.Exec(query, newUser.Name, newUser.PhoneNumber, newUser.Email, newUser.UserRole, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Update the address related to the user
	addressQuery := "UPDATE Address SET "
	params := []interface{}{}

	if newUser.Address.Country != "" {
		addressQuery += "Country = ?, "
		params = append(params, newUser.Address.Country)
	}
	if newUser.Address.City != "" {
		addressQuery += "City = ?, "
		params = append(params, newUser.Address.City)
	}
	if newUser.Address.State != "" {
		addressQuery += "State = ?, "
		params = append(params, newUser.Address.State)
	}
	if newUser.Address.Street != "" {
		addressQuery += "Street = ?, "
		params = append(params, newUser.Address.Street)
	}
	if newUser.Address.PostalCode != "" {
		addressQuery += "PostalCode = ?, "
		params = append(params, newUser.Address.PostalCode)
	}
	if newUser.Address.AdditionalNumber != "" {
		addressQuery += "AdditionalNumber = ?, "
		params = append(params, newUser.Address.AdditionalNumber)
	}
	if newUser.Address.MapLocation != "" {
		addressQuery += "MapLocation = ?, "
		params = append(params, newUser.Address.MapLocation)
	}
	if newUser.Address.Latitude != "" {
		addressQuery += "Latitude = ?, "
		params = append(params, newUser.Address.Latitude)
	}
	if newUser.Address.Longitude != "" {
		addressQuery += "Longitude = ?, "
		params = append(params, newUser.Address.Longitude)
	}

	// Remove the trailing comma and space
	addressQuery = addressQuery[:len(addressQuery)-2]

	addressQuery += " WHERE AddressID = ?"
	params = append(params, oldUser.AddressID)

	_, err = UserHandler.db.Exec(addressQuery, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to update address",
			Data:    err.Error(),
		})
		return
	}

	UserHandler.LoadUsersIntoCache()

	response := Response{
		Status:  "success",
		Message: "User and address updated successfully",
	}
	c.JSON(http.StatusOK, response)
}

func (UserHandler *UserHandler) DeleteUserHandler(c *gin.Context) {
	userID := c.Param("id")
	UserHandler.LoadUsersIntoCache()
	query := `DELETE FROM User WHERE UserID = ?`
	_, err := UserHandler.db.Exec(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	response := Response{
		Status:  "success",
		Message: "User deleted successfully",
		Data:    UserHandler.cache[userID],
	}
	UserHandler.LoadUsersIntoCache()
	c.JSON(http.StatusOK, response)
}

func (UserHandler *UserHandler) LoginHandler(c *gin.Context) {
	// Parse and decode the request body into a new 'User' struct
	var user Entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get the existing user details from the database
	var createTime []byte
	var existingUser Entities.User
	query := `SELECT * FROM User WHERE Email = ?`
	err := UserHandler.db.QueryRow(query, user.Email).Scan(&existingUser.UserID, &existingUser.AddressID, &existingUser.Name, &existingUser.PhoneNumber, &existingUser.Email, &existingUser.Password, &createTime, &existingUser.UserRole)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the user does not exist, send an appropriate response message
			response := Response{
				Status:  "error",
				Message: "User not found",
			}
			c.JSON(http.StatusUnauthorized, response)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}
	existingUser.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
	// Compare the supplied password with the stored password
	if user.Password != existingUser.Password {
		// If the password does not match, send an appropriate response message
		response := Response{
			Status:  "error",
			Message: "Invalid password",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// If the password matches, send a success response
	response := Response{
		Status:  "success",
		Message: "Logged in successfully",
		Data:    existingUser,
	}
	c.JSON(http.StatusOK, response)
}

//Completed by: Yousef Almutairi
