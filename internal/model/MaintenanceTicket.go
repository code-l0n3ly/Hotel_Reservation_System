package model

import (
	"errors"
	"time"
)

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

func (m *MaintenanceTicket) Validate() error {
	if m.TicketID == "" {
		return errors.New("TicketID is required")
	}
	if m.MaintenancePresenterID == "" {
		return errors.New("MaintenancePresenterID is required")
	}
	if m.TenantID == "" {
		return errors.New("TenantID is required")
	}
	if m.PropertyID == "" {
		return errors.New("PropertyID is required")
	}
	if m.Description == "" {
		return errors.New("description is required")
	}
	if m.UrgencyLevel == "" {
		return errors.New("UrgencyLevel is required")
	}
	if m.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

func (m *MaintenanceTicket) IsUrgent() bool {
	return m.UrgencyLevel == "high"
}

func (m *MaintenanceTicket) IsOpen() bool {
	return m.Status == "open"
}

func (m *MaintenanceTicket) IsClosed() bool {
	return m.Status == "closed"
}
