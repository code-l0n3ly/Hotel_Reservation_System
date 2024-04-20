package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterFinancialTransactionRoutes(router *mux.Router, FinancialTransactionHandler *handler.FinancialTransactionHandler) {
	router.HandleFunc("/financialTransaction/create", FinancialTransactionHandler.CreateTransaction).Methods("POST")
	router.HandleFunc("/financialTransaction/{id}", FinancialTransactionHandler.GetTransaction).Methods("GET")
	router.HandleFunc("/financialTransaction/{id}", FinancialTransactionHandler.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/financialTransaction/{id}", FinancialTransactionHandler.DeleteTransaction).Methods("DELETE")
}
