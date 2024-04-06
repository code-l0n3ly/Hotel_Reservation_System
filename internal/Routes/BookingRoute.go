package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterBookingRoutes(router *mux.Router, BookingHandler *handler.BookingHandler) {
	router.HandleFunc("/booking/create", BookingHandler.CreateBookingHandler).Methods("POST")
	router.HandleFunc("/booking/{id}", BookingHandler.GetBookingHandler).Methods("GET")
	router.HandleFunc("/booking/{id}", BookingHandler.UpdateBookingHandler).Methods("PUT")
	router.HandleFunc("/booking/{id}", BookingHandler.DeleteBookingHandler).Methods("DELETE")
}
