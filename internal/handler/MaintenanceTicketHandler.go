package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gorilla/mux"
)

type MaintenanceTicketHandler struct {
	db                *sql.DB
	TicketIdReference int64
	cache             map[string]Entities.MaintenanceTicket // Cache to hold tickets in memory
}

func NewMaintenanceTicketHandler(db *sql.DB) *MaintenanceTicketHandler {
	return &MaintenanceTicketHandler{
		db:                db,
		TicketIdReference: 0,
		cache:             make(map[string]Entities.MaintenanceTicket),
	}
}

func (handler *MaintenanceTicketHandler) GenerateUniqueTicketID() string {
	handler.TicketIdReference++
	return fmt.Sprintf("%d", handler.TicketIdReference)
}

func (handler *MaintenanceTicketHandler) SetHighestTicketID() {
	highestID := int64(0)
	for _, ticket := range handler.cache {
		ticketID, err := strconv.ParseInt(ticket.TicketID, 10, 64)
		if err != nil {
			continue // Skip if the TicketID is not a valid integer
		}
		if ticketID > highestID {
			highestID = ticketID
		}
	}
	handler.TicketIdReference = highestID
}

func (handler *MaintenanceTicketHandler) LoadTickets() error {
	rows, err := handler.db.Query(`SELECT TicketID, MaintenancePresenterID, TenantID, PropertyID, Description, UrgencyLevel, CreateTime, Status FROM MaintenanceTicket`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ticket Entities.MaintenanceTicket
		if err := rows.Scan(&ticket.TicketID, &ticket.MaintenancePresenterID, &ticket.TenantID, &ticket.PropertyID, &ticket.Description, &ticket.UrgencyLevel, &ticket.CreateTime, &ticket.Status); err != nil {
			return err
		}
		handler.cache[ticket.TicketID] = ticket
	}
	handler.SetHighestTicketID()
	return rows.Err()
}

func (handler *MaintenanceTicketHandler) CreateMaintenanceTicket(w http.ResponseWriter, r *http.Request) {
	var ticket Entities.MaintenanceTicket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	handler.LoadTickets()

	query := `INSERT INTO MaintenanceTicket (TicketID, MaintenancePresenterID, TenantID, PropertyID, Description, UrgencyLevel, CreateTime, Status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = handler.db.Exec(query, ticket.TicketID, ticket.MaintenancePresenterID, ticket.TenantID, ticket.PropertyID, ticket.Description, ticket.UrgencyLevel, time.Now(), ticket.Status)
	if err != nil {
		http.Error(w, "Failed to create maintenance ticket", http.StatusInternalServerError)
		return
	}
	handler.LoadTickets()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(handler.cache[ticket.TicketID]) // Respond with the created ticket object
}

func (handler *MaintenanceTicketHandler) GetMaintenanceTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ticketID := params["id"]

	var ticket Entities.MaintenanceTicket
	query := `SELECT TicketID, MaintenancePresenterID, TenantID, PropertyID, Description, UrgencyLevel, CreateTime, Status FROM MaintenanceTicket WHERE TicketID = ?`
	err := handler.db.QueryRow(query, ticketID).Scan(&ticket.TicketID, &ticket.MaintenancePresenterID, &ticket.TenantID, &ticket.PropertyID, &ticket.Description, &ticket.UrgencyLevel, &ticket.CreateTime, &ticket.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve maintenance ticket", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

func (handler *MaintenanceTicketHandler) UpdateMaintenanceTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ticketID := params["id"]
	handler.LoadTickets()
	var ticket Entities.MaintenanceTicket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE MaintenanceTicket SET MaintenancePresenterID = ?, TenantID = ?, PropertyID = ?, Description = ?, UrgencyLevel = ?, Status = ? WHERE TicketID = ?`
	_, err = handler.db.Exec(query, ticket.MaintenancePresenterID, ticket.TenantID, ticket.PropertyID, ticket.Description, ticket.UrgencyLevel, ticket.Status, ticketID)
	if err != nil {
		http.Error(w, "Failed to update maintenance ticket", http.StatusInternalServerError)
		return
	}
	handler.LoadTickets()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Maintenance ticket updated successfully")
}

func (handler *MaintenanceTicketHandler) DeleteMaintenanceTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ticketID := params["id"]
	handler.LoadTickets()
	query := `DELETE FROM MaintenanceTicket WHERE TicketID = ?`
	_, err := handler.db.Exec(query, ticketID)
	if err != nil {
		http.Error(w, "Failed to delete maintenance ticket", http.StatusInternalServerError)
		return
	}
	handler.LoadTickets()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Maintenance ticket deleted successfully")
}
