package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterMaintenanceTicketRoutes(router *mux.Router, MaintenanceTicketHandler *handler.MaintenanceTicketHandler) {
	router.HandleFunc("/maintenanceTicket/create", MaintenanceTicketHandler.CreateMaintenanceTicket).Methods("POST")
	router.HandleFunc("/maintenanceTicket/{id}", MaintenanceTicketHandler.GetMaintenanceTicket).Methods("GET")
	router.HandleFunc("/maintenanceTicket/{id}", MaintenanceTicketHandler.UpdateMaintenanceTicket).Methods("PUT")
	router.HandleFunc("/maintenanceTicket/{id}", MaintenanceTicketHandler.DeleteMaintenanceTicket).Methods("DELETE")
}
