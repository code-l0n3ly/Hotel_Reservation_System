package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterPropertyRoutes(router *mux.Router, PropertyHandler *handler.PropertyHandler) {
	router.HandleFunc("/property/create", PropertyHandler.CreateProperty).Methods("POST")
	router.HandleFunc("/property/{id}", PropertyHandler.GetProperty).Methods("GET")
	router.HandleFunc("/property/{id}", PropertyHandler.UpdateProperty).Methods("PUT")
	router.HandleFunc("/property/{id}", PropertyHandler.DeleteProperty).Methods("DELETE")
}
