package Handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	db    *sql.DB
	cache map[string]Entities.Booking // Cache to hold bookings in memory
}

func NewBookingHandler(db *sql.DB) *BookingHandler {
	return &BookingHandler{
		db:    db,
		cache: make(map[string]Entities.Booking),
	}
}

func (BookingHandler *BookingHandler) LoadBookings() error {
	rows, err := BookingHandler.db.Query(`SELECT BookingID, UnitID, UserID, EndDate, CreateTime, StartDate, Summary FROM Booking`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var StartDate []byte
		var EndDate []byte
		var booking Entities.Booking
		if err := rows.Scan(&booking.BookingID, &booking.UnitID, &booking.UserID, &EndDate, &createTime, &StartDate, &booking.Summary); err != nil {
			fmt.Println(err.Error())
			return err
		}
		booking.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		booking.StartDate, _ = time.Parse("2006-01-02 15:04:05", string(StartDate))
		booking.EndDate, _ = time.Parse("2006-01-02 15:04:05", string(EndDate))
		fmt.Println(booking)
		BookingHandler.cache[booking.BookingID] = booking
	}
	return rows.Err()
}

// Method takes date as an argument and check in the cache if there is an active booking in that date returns true otherwise false
func (BookingHandler *BookingHandler) CheckActiveBooking(UnitId string, date time.Time) bool {
	for _, booking := range BookingHandler.cache {
		if booking.UnitID == UnitId && booking.StartDate.Before(date) || booking.EndDate.After(date) || booking.StartDate.Equal(date) || booking.EndDate.Equal(date) {
			return true
		} else {
			return false
		}
	}
	return false
}
func (BookingHandler *BookingHandler) CreateBooking(c *gin.Context) {
	var booking Entities.Booking
	BookingHandler.LoadBookings()

	err := c.BindJSON(&booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if BookingHandler.CheckActiveBooking(booking.UnitID, booking.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "There is an active booking in this date"})
		return
	}
	query := `INSERT INTO Booking (UnitID, UserID, EndDate, StartDate, Summary) VALUES (?, ?, ?, ?, ?)`
	result, err := BookingHandler.db.Exec(query, booking.UnitID, booking.UserID, booking.EndDate, booking.StartDate, booking.Summary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create booking" + err.Error()})
		return
	}
	id, _ := result.LastInsertId()
	booking.BookingID = strconv.FormatInt(id, 10)
	BookingHandler.LoadBookings()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Booking created successfully", "data": BookingHandler.cache[booking.BookingID]})
}

func (BookingHandler *BookingHandler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	BookingHandler.LoadBookings()

	booking, exists := BookingHandler.cache[bookingID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Booking retrieved successfully", "data": booking})
}

func (BookingHandler *BookingHandler) UpdateBooking(c *gin.Context) {
	bookingID := c.Param("id")
	BookingHandler.LoadBookings()

	var newInfoBooking Entities.Booking
	oldInfoBooking := BookingHandler.cache[bookingID]

	err := c.BindJSON(&newInfoBooking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !newInfoBooking.StartDate.IsZero() {
		oldInfoBooking.StartDate = newInfoBooking.StartDate
	}
	if !newInfoBooking.EndDate.IsZero() {
		oldInfoBooking.EndDate = newInfoBooking.EndDate
	}
	if newInfoBooking.Summary != "" {
		oldInfoBooking.Summary = newInfoBooking.Summary
	}

	query := `UPDATE Booking SET StartDate = ?, EndDate = ?, Summary = ? WHERE BookingID = ?`
	_, err = BookingHandler.db.Exec(query, oldInfoBooking.StartDate, oldInfoBooking.EndDate, oldInfoBooking.Summary, oldInfoBooking.BookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update booking" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Booking updated successfully", "Data": oldInfoBooking})
}

func (BookingHandler *BookingHandler) DeleteBooking(c *gin.Context) {
	bookingID := c.Param("id")
	BookingHandler.LoadBookings()

	query := `DELETE FROM Booking WHERE BookingID = ?`
	_, err := BookingHandler.db.Exec(query, bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to delete booking" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Booking deleted successfully", "data": BookingHandler.cache[bookingID]})
	BookingHandler.LoadBookings()
}
