package Handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	UserHandler.cache = make(map[string]Entities.User)
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
	if len(params) > 0 {
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
	}

	UserHandler.LoadUsersIntoCache()

	response := Response{
		Status:  "success",
		Message: "User and address updated successfully",
		Data:    UserHandler.cache[userID],
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

type Report struct {
	ReportID              string `json:"reportID"`
	UserID                string `json:"userID"`
	Type                  string `json:"type,omitempty"`
	CreateTime            string `json:"createTime,omitempty"`
	Properties            []Entities.Property
	Bookings              []Entities.Booking
	FinancialTransactions []Entities.FinancialTransaction
	TotalEarnings         int `json:"totalEarnings,omitempty"`
}

func (UserHandler *UserHandler) GetProperties(userID string) ([]Entities.Property, error) {
	query := `
    SELECT 
        p.PropertyID, 
        p.OwnerID, 
        p.Name, 
        p.CreateTime,
        a.AddressID, 
        a.Country, 
        a.City, 
        a.State, 
        a.Street, 
        a.PostalCode, 
        a.AdditionalNumber, 
        a.MapLocation, 
        a.Latitude, 
        a.Longitude
    FROM 
        Property p
    LEFT JOIN 
        Address a ON p.AddressID = a.AddressID
    WHERE 
        p.OwnerID = ?
    `
	rows, err := UserHandler.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []Entities.Property
	for rows.Next() {
		var createTime []byte
		var property Entities.Property
		var address Entities.Address
		if err := rows.Scan(&property.PropertyID, &property.OwnerID, &property.Name, &createTime, &address.AddressID, &address.Country, &address.City, &address.State, &address.Street, &address.PostalCode, &address.AdditionalNumber, &address.MapLocation, &address.Latitude, &address.Longitude); err != nil {
			return nil, err
		}
		property.CreateTime, err = time.Parse("2006-01-02 15:04:05", string(createTime))
		if err != nil {
			return nil, err
		}
		property.Address = address

		// Get units for the property
		units, err := UserHandler.GetUnits(property.PropertyID)
		if err != nil {
			return nil, err
		}
		property.Units = units

		properties = append(properties, property)
	}

	return properties, nil
}

func (UserHandler *UserHandler) GetReports(c *gin.Context) {
	UserHandler.LoadUsersIntoCache()
	userID := c.Param("id")
	user, exists := UserHandler.GetUserByID(userID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	Report, _ := UserHandler.GetReport(user.UserID)
	//fmt.Println(err.Error())
	c.JSON(http.StatusOK, Report)
}

func (UserHandler *UserHandler) GetUnits(propertyID string) ([]Entities.Unit, error) {
	query := `
    SELECT 
        u.UnitID, 
        u.PropertyID, 
        u.Name, 
        u.RentalPrice, 
        u.Description, 
        u.Rating, 
        u.StructuralProperties, 
        u.CreateTime,
        a.AddressID, 
        a.Country, 
        a.City, 
        a.State, 
        a.Street, 
        a.PostalCode, 
        a.AdditionalNumber, 
        a.MapLocation, 
        a.Latitude, 
        a.Longitude
    FROM 
        Unit u
    LEFT JOIN 
        Address a ON u.AddressID = a.AddressID
    WHERE 
        u.PropertyID = ?
    `
	rows, err := UserHandler.db.Query(query, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Entities.Unit
	for rows.Next() {
		var createTime []byte
		var unit Entities.Unit
		var address Entities.Address
		if err := rows.Scan(&unit.UnitID, &unit.PropertyID, &unit.Name, &unit.RentalPrice, &unit.Description, &unit.Rating, &unit.StructuralProperties, &createTime, &address.AddressID, &address.Country, &address.City, &address.State, &address.Street, &address.PostalCode, &address.AdditionalNumber, &address.MapLocation, &address.Latitude, &address.Longitude); err != nil {
			return nil, err
		}
		unit.CreateTime, err = time.Parse("2006-01-02 15:04:05", string(createTime))
		if err != nil {
			return nil, err
		}
		unit.Address = address
		units = append(units, unit)
	}

	return units, nil
}

// Create a report and return it for a userid
func (UserHandler *UserHandler) GetReport(userID string) (Report, error) {
	// Get the properties for the user
	properties, err := UserHandler.GetProperties(userID)
	if err != nil {
		return Report{}, err
	}

	// Get the bookings for the user
	var bookings []Entities.Booking

	for _, property := range properties {
		for _, unit := range property.Units {
			unitBookings, err := UserHandler.GetBookings(unit.UnitID)
			if err != nil {
				continue
			}
			bookings = append(bookings, unitBookings...)
		}
	}
	var FinancialTransactions []Entities.FinancialTransaction
	for _, booking := range bookings {
		transactions, err := UserHandler.GetFinancialTransactions(booking.BookingID)
		if err != nil {
			continue
		}
		FinancialTransactions = append(FinancialTransactions, transactions...)
	}

	// Calculate the total earnings for the user
	var totalEarnings int
	for _, Transaction := range FinancialTransactions {
		totalEarnings += Transaction.Amount
	}

	report := Report{
		ReportID:              uuid.New().String(),
		UserID:                userID,
		Type:                  "Report",
		CreateTime:            time.Now().Format("2006-01-02 15:04:05"),
		Properties:            properties,
		Bookings:              bookings,
		FinancialTransactions: FinancialTransactions,
		TotalEarnings:         totalEarnings,
	}

	return report, nil
}

func (UserHandler *UserHandler) GetBookings(UnitID string) ([]Entities.Booking, error) {
	query := `SELECT * FROM Booking WHERE UnitID = ?`
	rows, err := UserHandler.db.Query(query, UnitID)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()
	var bookings []Entities.Booking
	for rows.Next() {
		var booking Entities.Booking
		var createTime []byte
		var StartDate []byte
		var EndDate []byte
		if err := rows.Scan(&booking.BookingID, &booking.UnitID, &booking.UserID, &EndDate, &createTime, &StartDate, &booking.Summary); err != nil {
			return nil, err
		}
		booking.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		booking.StartDate, _ = time.Parse("2006-01-02 15:04:05", string(StartDate))
		booking.EndDate, _ = time.Parse("2006-01-02 15:04:05", string(EndDate))
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

// Gets all the FinancialTransactions for the bookings for the units of a property
func (UserHandler *UserHandler) GetFinancialTransactions(BookingID string) ([]Entities.FinancialTransaction, error) {
	var transactions []Entities.FinancialTransaction

	query := `SELECT TransactionID, UserID, BookingID, PaymentMethod, Amount, CreateTime FROM FinancialTransaction WHERE BookingID = ?`
	rows, err := UserHandler.db.Query(query, BookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Entities.FinancialTransaction
		var createTime []byte
		if err := rows.Scan(&transaction.TransactionID, &transaction.UserID, &transaction.BookingID, &transaction.PaymentMethod, &transaction.Amount, &createTime); err != nil {
			return nil, err
		}
		transaction.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
