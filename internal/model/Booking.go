package model

import (
	"errors"
	"time"
)

// Booking represents the 'Booking' table in your database.
type Booking struct {
	BookingID  string    `json:"bookingID"`
	UnitID     string    `json:"unitID"`
	UserID     string    `json:"userID"`
	EndDate    time.Time `json:"endDate"`
	CreateTime time.Time `json:"createTime"`
	StartDate  time.Time `json:"startDate"`
	Summary    string    `json:"summary"` // Assuming JSON data as a string; adjust according to your needs
}

func (b *Booking) Validate() error {
	if b.BookingID == "" {
		return errors.New("BookingID is required")
	}
	if b.UnitID == "" {
		return errors.New("UnitID is required")
	}
	if b.UserID == "" {
		return errors.New("UserID is required")
	}
	if b.EndDate.IsZero() {
		return errors.New("EndDate is required")
	}
	if b.StartDate.IsZero() {
		return errors.New("StartDate is required")
	}
	if b.StartDate.After(b.EndDate) {
		return errors.New("StartDate cannot be after EndDate")
	}
	if b.Summary == "" {
		return errors.New("summary is required")
	}
	return nil
}

func (b *Booking) IsPastBooking() bool {
	return time.Now().After(b.EndDate)
}

func (b *Booking) IsFutureBooking() bool {
	return time.Now().Before(b.StartDate)
}

func (b *Booking) IsActiveBooking() bool {
	now := time.Now()
	return now.After(b.StartDate) && now.Before(b.EndDate)
}
