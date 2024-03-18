package UserHandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "GraduationProject.com/m/internal/db"
	Entities "GraduationProject.com/m/internal/model"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (UserHandler *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user Entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO User (UserID, Name, Email, Password) VALUES (?, ?, ?, ?)`
	response, err := UserHandler.db.Exec(query, user.UserID, user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	resp, _ := response.LastInsertId()
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte((fmt.Sprintf("%d", resp)))) // Respond with the ID of the created user (assuming it's an auto-incremented ID
	//json.NewEncoder(w).Encode(user) // Respond with the created user object
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
