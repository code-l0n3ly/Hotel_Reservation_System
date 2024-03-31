package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes sets up the routes for the application
func RegisterUserRoutes(router *mux.Router, UserHandler *handler.UserHandler) {
	router.HandleFunc("/users", UserHandler.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/{id}", UserHandler.GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{id}", UserHandler.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", UserHandler.DeleteUserHandler).Methods("DELETE")
}
