package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterReportRoutes(router *mux.Router, ReportHandler *handler.ReportHandler) {
	router.HandleFunc("/report/create", ReportHandler.CreateReport).Methods("POST")
	router.HandleFunc("/report/{id}", ReportHandler.GetReport).Methods("GET")
	router.HandleFunc("/report/{id}", ReportHandler.UpdateReport).Methods("PUT")
	router.HandleFunc("/report/{id}", ReportHandler.DeleteReport).Methods("DELETE")
}
