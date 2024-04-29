package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "GraduationProject.com/m/internal/db"
	Entities "GraduationProject.com/m/internal/model"
	"github.com/gorilla/mux"
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

func (UserHandler *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	UserHandler.LoadUsersIntoCache()
	var user Entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	user.UserID = UserHandler.GenerateUniqueUserID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists in the cache
	user.UserID = strconv.Itoa(int(UserHandler.UserIdReference) + 1)
	_, exists := UserHandler.GetUserByID(user.UserID)
	if exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !user.IsEmailValid() {
		http.Error(w, "Email is not valid", http.StatusBadRequest)
		return
	} else if !user.IsPasswordStrong() {
		http.Error(w, "Password is not strong enough", http.StatusBadRequest)
		return
	}
	query := `INSERT INTO User (UserID, Name, Email, PhoneNumber Password, UserRole) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = UserHandler.db.Exec(query, user.UserID, user.Name, user.Email, user.PhoneNumber, user.Password, user.UserRole)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Add the new user to the cache
	UserHandler.LoadUsersIntoCache()
	w.WriteHeader(http.StatusCreated)
	response := Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    UserHandler.cache[user.UserID],
	}

	// Write the response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (UserHandler *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var user Entities.User
	query := `SELECT UserID, Name, Email, PhoneNumber, UserRole FROM User WHERE UserID = ?`
	err := UserHandler.db.QueryRow(query, userID).Scan(&user.UserID, &user.Name, &user.Email, &user.PhoneNumber, &user.UserRole)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	response := Response{
		Status:  "success",
		Message: "User retrieved successfully",
		Data:    user,
	}
	json.NewEncoder(w).Encode(response)
}

func (UserHandler *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(response)
}

func (UserHandler *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]
	UserHandler.LoadUsersIntoCache()
	var user Entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE User SET Name = ?, Email = ?, UserRole = ? WHERE UserID = ?`
	_, err = UserHandler.db.Exec(query, user.Name, user.Email, user.UserRole, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	UserHandler.LoadUsersIntoCache()
	w.WriteHeader(http.StatusOK)
	response := Response{
		Status:  "success",
		Message: "User updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (UserHandler *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]
	UserHandler.LoadUsersIntoCache()
	query := `DELETE FROM User WHERE UserID = ?`
	_, err := UserHandler.db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := Response{
		Status:  "success",
		Message: "User deleted successfully",
		Data:    UserHandler.cache[userID],
	}
	UserHandler.LoadUsersIntoCache()
	json.NewEncoder(w).Encode(response)
}

func (UserHandler *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new 'User' struct
	var user Entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	UserHandler.LoadUsersIntoCache()
	// Get the existing user details from the database
	var existingUser Entities.User
	query := `SELECT UserID, Name, Email, Password FROM User WHERE Email = ?`
	err = UserHandler.db.QueryRow(query, user.Email).Scan(&existingUser.UserID, &existingUser.Name, &existingUser.Email, &existingUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the user does not exist, send an appropriate response message
			response := Response{
				Status:  "error",
				Message: "User not found",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// If the password matches, send a success response
	response := Response{
		Status:  "success",
		Message: "Logged in successfully",
		Data:    existingUser,
	}
	json.NewEncoder(w).Encode(response)
}
