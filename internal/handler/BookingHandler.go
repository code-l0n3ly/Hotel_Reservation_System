package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gorilla/mux"
)

type BookingHandler struct {
	db                 *sql.DB
	BookingIdReference int64
	cache              map[string]Entities.Booking // Cache to hold bookings in memory
}

func NewBookingHandler(db *sql.DB) *BookingHandler {
	return &BookingHandler{
		db:                 db,
		BookingIdReference: 0,
		cache:              make(map[string]Entities.Booking),
	}
}

func (handler *BookingHandler) GenerateUniqueBookingID() string {
	handler.BookingIdReference++
	return fmt.Sprintf("%d", handler.BookingIdReference)
}

func (handler *BookingHandler) SetHighestBookingID() {
	highestID := int64(0)
	for _, booking := range handler.cache {
		bookingID, err := strconv.ParseInt(booking.BookingID, 10, 64)
		if err != nil {
			continue // Skip if the BookingID is not a valid integer
		}
		if bookingID > highestID {
			highestID = bookingID
		}
	}
	handler.BookingIdReference = highestID
}

func (handler *BookingHandler) LoadBookings() error {
	rows, err := handler.db.Query(`SELECT BookingID, UnitID, UserID, StartDate, EndDate, CreateTime, Summary FROM Booking`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var booking Entities.Booking
		if err := rows.Scan(&booking.BookingID, &booking.UnitID, &booking.UserID, &booking.StartDate, &booking.EndDate, &booking.CreateTime, &booking.Summary); err != nil {
			return err
		}
		handler.cache[booking.BookingID] = booking
	}
	handler.SetHighestBookingID()
	return rows.Err()
}

func (handler *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking Entities.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	handler.LoadBookings()

	query := `INSERT INTO Booking (BookingID, UnitID, UserID, StartDate, EndDate, CreateTime, Summary) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = handler.db.Exec(query, booking.BookingID, booking.UnitID, booking.UserID, booking.StartDate, booking.EndDate, booking.CreateTime, booking.Summary)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}
	handler.LoadBookings()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(handler.cache[booking.BookingID]) // Respond with the created booking object
}

func (handler *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID := params["id"]

	var booking Entities.Booking
	query := `SELECT BookingID, UnitID, UserID, StartDate, EndDate, CreateTime, Summary FROM Booking WHERE BookingID = ?`
	err := handler.db.QueryRow(query, bookingID).Scan(&booking.BookingID, &booking.UnitID, &booking.UserID, &booking.StartDate, &booking.EndDate, &booking.CreateTime, &booking.Summary)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve booking", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(booking)
}

func (handler *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID := params["id"]
	handler.LoadBookings()
	var booking Entities.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE Booking SET UnitID = ?, UserID = ?, StartDate = ?, EndDate = ?, CreateTime = ?, Summary = ? WHERE BookingID = ?`
	_, err = handler.db.Exec(query, booking.UnitID, booking.UserID, booking.StartDate, booking.EndDate, booking.CreateTime, booking.Summary, bookingID)
	if err != nil {
		http.Error(w, "Failed to update booking", http.StatusInternalServerError)
		return
	}
	handler.LoadBookings()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Booking updated successfully")
}

func (handler *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID := params["id"]
	handler.LoadBookings()
	query := `DELETE FROM Booking WHERE BookingID = ?`
	_, err := handler.db.Exec(query, bookingID)
	if err != nil {
		http.Error(w, "Failed to delete booking", http.StatusInternalServerError)
		return
	}
	handler.LoadBookings()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Booking deleted successfully")
}
