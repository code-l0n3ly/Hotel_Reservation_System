package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterMessageRoutes(router *mux.Router, MessageHandler *handler.MessageHandler) {
	router.HandleFunc("/message/create", MessageHandler.CreateMessage).Methods("POST")
	router.HandleFunc("/message/{id}", MessageHandler.GetMessage).Methods("GET")
	router.HandleFunc("/message/{id}", MessageHandler.UpdateMessage).Methods("PUT")
	router.HandleFunc("/message/{id}", MessageHandler.DeleteMessage).Methods("DELETE")
}
