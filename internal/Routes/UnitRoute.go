package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterUnitRoutes(router *mux.Router, UnitHandler *handler.UnitHandler) {
	router.HandleFunc("/units", UnitHandler.CreateUnit).Methods("POST")
	router.HandleFunc("/units/{id}", UnitHandler.GetUnit).Methods("GET")
	router.HandleFunc("/units/{id}", UnitHandler.UpdateUnit).Methods("PUT")
	router.HandleFunc("/units/{id}", UnitHandler.DeleteUnit).Methods("DELETE")
}
