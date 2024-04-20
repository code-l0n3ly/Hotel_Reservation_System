package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterBookingRoutes(router *mux.Router, BookingHandler *handler.BookingHandler) {
	router.HandleFunc("/booking/create", BookingHandler.CreateBooking).Methods("POST")
	router.HandleFunc("/booking/{id}", BookingHandler.GetBooking).Methods("GET")
	router.HandleFunc("/booking/{id}", BookingHandler.UpdateBooking).Methods("PUT")
	router.HandleFunc("/booking/{id}", BookingHandler.DeleteBooking).Methods("DELETE")
}
