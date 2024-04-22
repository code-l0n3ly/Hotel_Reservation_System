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

type UnitHandler struct {
	db              *sql.DB
	UnitIdReference int64
	cache           map[string]Entities.Unit // Cache to hold users in memory
}

func NewUnitHandler(db *sql.DB) *UnitHandler {
	return &UnitHandler{
		db:              db,
		UnitIdReference: 0,
		cache:           make(map[string]Entities.Unit),
	}
}

func (UnitHandler *UnitHandler) GenerateUniqueUnitID() string {
	UnitHandler.UnitIdReference++
	return fmt.Sprintf("%d", UnitHandler.UnitIdReference)
}

func (UnitHandler *UnitHandler) SetHighestUnitID() {
	highestID := int64(0)
	for _, unit := range UnitHandler.cache {
		unitID, err := strconv.ParseInt(unit.UnitID, 10, 64)
		if err != nil {
			continue // Skip if the UnitID is not a valid integer
		}
		if unitID > highestID {
			highestID = unitID
		}
	}
	UnitHandler.UnitIdReference = highestID
}

func (UnitHandler *UnitHandler) LoadUnits() error {
	rows, err := UnitHandler.db.Query(`SELECT UnitID, PropertyID, RentalPrice, OccupancyStatus, StructuralProperties, CreateTime FROM Unit`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var unit Entities.Unit
		if err := rows.Scan(&unit.UnitID, &unit.PropertyID, &unit.RentalPrice, &unit.OccupancyStatus, &unit.StructuralProperties, &unit.CreateTime); err != nil {
			return err
		}
		fmt.Println(unit)
		UnitHandler.cache[unit.UnitID] = unit
	}
	UnitHandler.SetHighestUnitID()
	return rows.Err()
}

func (UnitHandler *UnitHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var unit Entities.Unit
	UnitHandler.LoadUnits()
	unit.UnitID = UnitHandler.GenerateUniqueUnitID()
	err := json.NewDecoder(r.Body).Decode(&unit)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = unit.Validate()
	if err != nil {
		response := Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	query := `INSERT INTO Unit (UnitID, Images, Description, Rating, PropertyID, RentalPrice, OccupancyStatus, StructuralProperties, CreateTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	tx, err := UnitHandler.db.Begin()
	if err != nil {
		http.Error(w, "Failed to create unit", http.StatusInternalServerError)
	}

	_, err = tx.Exec(query, unit.UnitID, unit.Description, unit.Rating, unit.PropertyID, unit.RentalPrice, unit.OccupancyStatus, unit.StructuralProperties, time.Now())
	if err != nil {
		http.Error(w, "Failed to create unit", http.StatusInternalServerError)
	}

	for _, image := range unit.Images {
		_, err = tx.Exec(`INSERT INTO Images (UnitID, Image) VALUES (?, ?)`, unit.UnitID, image)
		if err != nil {
			http.Error(w, "Failed to create unit", http.StatusInternalServerError)
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to create unit", http.StatusInternalServerError)
	}
	UnitHandler.LoadUnits()
	w.WriteHeader(http.StatusCreated)
	response := Response{
		Status:  "success",
		Message: "Unit created successfully",
		Data:    UnitHandler.cache[unit.UnitID],
	}
	json.NewEncoder(w).Encode(response)
}

func (UnitHandler *UnitHandler) GetUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	unitID := params["id"]
	UnitHandler.LoadUnits()
	var unit Entities.Unit
	query := `SELECT UnitID, PropertyID, RentalPrice, OccupancyStatus, StructuralProperties, CreateTime FROM Unit WHERE UnitID = ?`
	err := UnitHandler.db.QueryRow(query, unitID).Scan(&unit.UnitID, &unit.PropertyID, &unit.RentalPrice, &unit.OccupancyStatus, &unit.StructuralProperties, &unit.CreateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve unit", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(unit)
}

func (UnitHandler *UnitHandler) UpdateUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	unitID := params["id"]
	UnitHandler.LoadUnits()
	var unit Entities.Unit
	err := json.NewDecoder(r.Body).Decode(&unit)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = unit.Validate()
	if err != nil {
		response := Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	query := `UPDATE Unit SET PropertyID = ?, RentalPrice = ?, OccupancyStatus = ?, StructuralProperties = ?, CreateTime = ? WHERE UnitID = ?`
	_, err = UnitHandler.db.Exec(query, unit.PropertyID, unit.RentalPrice, unit.OccupancyStatus, unit.StructuralProperties, unit.CreateTime, unitID)
	if err != nil {
		response := Response{
			Status:  "Failed to update unit",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Unit updated successfully")
}

func (UnitHandler *UnitHandler) DeleteUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	unitID := params["id"]
	UnitHandler.LoadUnits()
	query := `DELETE FROM Unit WHERE UnitID = ?`
	_, err := UnitHandler.db.Exec(query, unitID)
	if err != nil {
		response := Response{
			Status:  "Failed to delete unit",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	UnitHandler.LoadUnits()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Unit deleted successfully")
}

func (UnitHandler *UnitHandler) GetUnits(w http.ResponseWriter, r *http.Request) {
	UnitHandler.LoadUnits()
	var units []Entities.Unit
	for _, unit := range UnitHandler.cache {
		units = append(units, unit)
	}

	json.NewEncoder(w).Encode(units)
}
