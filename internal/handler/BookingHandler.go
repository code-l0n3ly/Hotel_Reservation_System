package Handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gorilla/mux"
)

type BookingHandler struct {
	db *sql.DB
}

func NewBookingHandler(db *sql.DB) *BookingHandler {
	return &BookingHandler{db: db}
}

func (handler *BookingHandler) CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	var booking Entities.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if there's already a booking for the same unit at the same time
	var existingBookingID string
	err = handler.db.QueryRow(`SELECT BookingID FROM Booking WHERE UnitID = ? AND ((StartDate <= ? AND EndDate >= ?) OR (StartDate <= ? AND EndDate >= ?))`, booking.UnitID, booking.StartDate, booking.StartDate, booking.EndDate, booking.EndDate).Scan(&existingBookingID)
	if err != sql.ErrNoRows {
		if err != nil {
			http.Error(w, "Failed to check existing bookings", http.StatusInternalServerError)
			return
		}
		response := Response{
			Status:  "error",
			Message: "There's already a booking for the same unit at the same time",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	query := `INSERT INTO Booking (BookingID, UserID, UnitID, StartDate, EndDate, CreateTime, Summary) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = handler.db.Exec(query, booking.BookingID, booking.UserID, booking.UnitID, booking.StartDate, booking.EndDate, booking.CreateTime, booking.Summary)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Booking created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (handler *BookingHandler) GetBookingHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID := params["id"]

	var booking Entities.Booking
	err := handler.db.QueryRow(`SELECT BookingID, UserID, UnitID, StartDate, EndDate, CreateTime, Summary FROM Booking WHERE BookingID = ?`, bookingID).Scan(&booking.BookingID, &booking.UserID, &booking.UnitID, &booking.StartDate, &booking.EndDate, &booking.CreateTime, &booking.Summary)
	if err != nil {
		http.Error(w, "Failed to get booking", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(booking)
}

// Implement other handlers (UpdateBookingHandler, DeleteBookingHandler, etc.) as necessary
func (handler *BookingHandler) UpdateBookingHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID := params["id"]

	var booking Entities.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE Booking SET UserID = ?, UnitID = ?, StartDate = ?, EndDate = ?, Summary = ? WHERE BookingID = ?`
	_, err = handler.db.Exec(query, booking.UserID, booking.UnitID, booking.StartDate, booking.EndDate, booking.Summary, bookingID)
	if err != nil {
		http.Error(w, "Failed to update booking", http.StatusInternalServerError)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Booking updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (handler *BookingHandler) DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID := params["id"]

	query := `DELETE FROM Booking WHERE BookingID = ?`
	_, err := handler.db.Exec(query, bookingID)
	if err != nil {
		http.Error(w, "Failed to delete booking", http.StatusInternalServerError)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Booking deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
