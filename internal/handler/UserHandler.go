package Handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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

func (UserHandler *UserHandler) GenerateUniqueUserID() string {
	UserHandler.UserIdReference++
	return fmt.Sprintf("%d", UserHandler.UserIdReference)
}

func (UserHandler *UserHandler) SetHighestUserID() {
	highestID := int64(0)
	for _, user := range UserHandler.cache {
		userID, err := strconv.ParseInt(user.UserID, 10, 64)
		if err != nil {
			continue // Skip if the UserID is not a valid integer
		}
		if userID > highestID {
			highestID = userID
		}
	}
	UserHandler.UserIdReference = highestID
}

func (UserHandler *UserHandler) LoadUsersIntoCache() error {

	rows, err := UserHandler.db.Query(`SELECT * FROM User`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var user Entities.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.PhoneNumber, &user.Password, createTime, &user.UserRole); err != nil {
			return err
		}
		user.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		fmt.Println(user)
		UserHandler.cache[user.UserID] = user
	}
	UserHandler.SetHighestUserID()
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

	user.UserID = UserHandler.GenerateUniqueUserID()

	// Check if user already exists in the cache
	user.UserID = strconv.Itoa(int(UserHandler.UserIdReference) + 1)
	_, exists := UserHandler.GetUserByID(user.UserID)
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if !user.IsEmailValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is not valid"})
		return
	} else if !user.IsPasswordStrong() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong enough"})
		return
	}
	query := `INSERT INTO User (UserID, Name, Email, PhoneNumber Password, UserRole) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := UserHandler.db.Exec(query, user.UserID, user.Name, user.Email, user.PhoneNumber, user.Password, user.UserRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Add the new user to the cache
	UserHandler.LoadUsersIntoCache()

	response := Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    UserHandler.cache[user.UserID],
	}

	c.JSON(http.StatusCreated, response)
}

func (UserHandler *UserHandler) GetUserHandler(c *gin.Context) {
	userID := c.Param("id")

	var user Entities.User
	query := `SELECT UserID, Name, Email, PhoneNumber, UserRole FROM User WHERE UserID = ?`
	err := UserHandler.db.QueryRow(query, userID).Scan(&user.UserID, &user.Name, &user.Email, &user.PhoneNumber, &user.UserRole)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
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
	var user Entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE User SET Name = ?, Email = ?, UserRole = ? WHERE UserID = ?`
	_, err := UserHandler.db.Exec(query, user.Name, user.Email, user.UserRole, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	UserHandler.LoadUsersIntoCache()

	response := Response{
		Status:  "success",
		Message: "User updated successfully",
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
	UserHandler.LoadUsersIntoCache()
	// Get the existing user details from the database
	var existingUser Entities.User
	query := `SELECT UserID, Name, Email, Password FROM User WHERE Email = ?`
	err := UserHandler.db.QueryRow(query, user.Email).Scan(&existingUser.UserID, &existingUser.Name, &existingUser.Email, &existingUser.Password)
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
