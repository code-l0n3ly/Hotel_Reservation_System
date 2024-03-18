package model

import (
	"database/sql"
	"time"
)

// User represents the 'User' table in your database.
type User struct {
	UserID     string    `json:"userID"`
	Name       string    `json:"name"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UserRole   string    `json:"userRole"`
}

// Report represents the 'Report' table in your database.
type Report struct {
	ReportID   string         `json:"reportID"`
	UserID     string         `json:"userID"`
	Type       sql.NullString `json:"type,omitempty"`
	CreateTime time.Time      `json:"createTime,omitempty"`
	Data       string         `json:"data"` // Assuming JSON data as a string; adjust according to your needs
}

// Property represents the 'Property' table in your database.
type Property struct {
	PropertyID  string         `json:"propertyID"`
	Name        sql.NullString `json:"name,omitempty"`
	Address     sql.NullString `json:"address,omitempty"`
	CreateTime  time.Time      `json:"createTime"`
	Type        string         `json:"type"`
	Photos      sql.NullString `json:"photos,omitempty"` // Consider changing to []string if storing multiple photos
	Description sql.NullString `json:"description,omitempty"`
	Rules       sql.NullString `json:"rules,omitempty"` // Assuming JSON data as a string; adjust according to your needs
}

// Unit represents the 'Unit' table in your database.
type Unit struct {
	UnitID               string    `json:"unitID"`
	PropertyID           string    `json:"propertyID"`
	RentalPrice          int       `json:"rentalPrice"`
	OccupancyStatus      string    `json:"occupancyStatus"`
	StructuralProperties string    `json:"structuralProperties"` // Assuming JSON data as a string; adjust according to your needs
	CreateTime           time.Time `json:"createTime"`
}

// Review represents the 'Review' table in your database.
type Review struct {
	ReviewID   string         `json:"reviewID"`
	UserID     string         `json:"userID"`
	UnitID     string         `json:"unitID"`
	Rating     int            `json:"rating"`
	Comment    sql.NullString `json:"comment,omitempty"`
	CreateTime time.Time      `json:"createTime"`
}

// Message represents the 'Message' table in your database.
type Message struct {
	MessageID  string    `json:"messageID"`
	Content    string    `json:"content"` // Assuming JSON data as a string; adjust according to your needs
	CreateTime time.Time `json:"createTime"`
	ReceiverID string    `json:"receiverID"`
	SenderID   string    `json:"senderID"`
}

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

// MaintenanceTicket represents the 'MaintenanceTicket' table in your database.
type MaintenanceTicket struct {
	TicketID               string    `json:"ticketID"`
	MaintenancePresenterID string    `json:"maintenancePresenterID"`
	TenantID               string    `json:"tenantID"`
	PropertyID             string    `json:"propertyID"`
	Description            string    `json:"description"`
	UrgencyLevel           string    `json:"urgencyLevel"`
	CreateTime             time.Time `json:"createTime"`
	Status                 string    `json:"status"`
}

// FinancialTransaction represents the 'FinancialTransaction' table in your database.
type FinancialTransaction struct {
	TransactionID string    `json:"transactionID"`
	UserID        string    `json:"userID"`
	UnitID        string    `json:"unitID"`
	PaymentMethod string    `json:"paymentMethod"`
	Amount        int       `json:"amount"`
	CreateTime    time.Time `json:"createTime"`
}
