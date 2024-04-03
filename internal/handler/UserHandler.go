package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	rows, err := UserHandler.db.Query(`SELECT UserID, Name, Email, UserRole FROM User`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var user Entities.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.CreateTime, &user.UserRole); err != nil {
			return err
		}
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
	_, exists := UserHandler.GetUserByID(user.UserID)
	if exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO User (UserID, Name, Email, Password) VALUES (?, ?, ?, ?)`
	_, err = UserHandler.db.Exec(query, user.UserID, user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Add the new user to the cache
	UserHandler.LoadUsersIntoCache()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserHandler.cache[user.UserID]) // Respond with the created user object
}

func (UserHandler *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var user Entities.User
	query := `SELECT UserID, Name, Email, UserRole FROM User WHERE UserID = ?`
	err := UserHandler.db.QueryRow(query, userID).Scan(&user.UserID, &user.Name, &user.Email, &user.UserRole)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (UserHandler *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	err := UserHandler.LoadUsersIntoCache()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	var users []Entities.User
	for _, user := range UserHandler.cache {
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func (UserHandler *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User updated successfully")
}

func (UserHandler *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	query := `DELETE FROM User WHERE UserID = ?`
	_, err := UserHandler.db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User deleted successfully")
}
