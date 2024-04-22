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

type PropertyHandler struct {
	db                  *sql.DB
	PropertyIdReference int64
	cache               map[string]Entities.Property // Cache to hold properties in memory
}

func NewPropertyHandler(db *sql.DB) *PropertyHandler {
	return &PropertyHandler{
		db:                  db,
		PropertyIdReference: 0,
		cache:               make(map[string]Entities.Property),
	}
}

func (PropertyHandler *PropertyHandler) GenerateUniquePropertyID() string {
	PropertyHandler.PropertyIdReference++
	return fmt.Sprintf("%d", PropertyHandler.PropertyIdReference)
}

func (PropertyHandler *PropertyHandler) SetHighestPropertyID() {
	highestID := int64(0)
	for _, property := range PropertyHandler.cache {
		propertyID, err := strconv.ParseInt(property.PropertyID, 10, 64)
		if err != nil {
			continue // Skip if the PropertyID is not a valid integer
		}
		if propertyID > highestID {
			highestID = propertyID
		}
	}
	PropertyHandler.PropertyIdReference = highestID
}

func (PropertyHandler *PropertyHandler) LoadProperties() error {
	rows, err := PropertyHandler.db.Query(`SELECT PropertyID, Name, Address, CreateTime, Type, Photos, Description, Rules FROM Property`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var property Entities.Property
		if err := rows.Scan(&property.PropertyID, &property.Name, &property.Address, &property.CreateTime, &property.Type, &property.Photos, &property.Description, &property.Rules); err != nil {
			return err
		}
		PropertyHandler.cache[property.PropertyID] = property
	}
	PropertyHandler.SetHighestPropertyID()
	return rows.Err()
}

func (PropertyHandler *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var property Entities.Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	PropertyHandler.LoadProperties()

	query := `INSERT INTO Property (PropertyID, Name, Address, Type, Description, Rules) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = PropertyHandler.db.Exec(query, property.PropertyID, property.Name, property.Address, property.Type, property.Description, property.Rules)
	if err != nil {
		http.Error(w, "Failed to create property"+err.Error(), http.StatusInternalServerError)
		return
	}
	PropertyHandler.LoadProperties()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PropertyHandler.cache[property.PropertyID]) // Respond with the created property object
}

func (PropertyHandler *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	propertyID := params["id"]

	var property Entities.Property
	query := `SELECT PropertyID, Name, Address, CreateTime, Type, Photos, Description, Rules FROM Property WHERE PropertyID = ?`
	err := PropertyHandler.db.QueryRow(query, propertyID).Scan(&property.PropertyID, &property.Name, &property.Address, &property.CreateTime, &property.Type, &property.Photos, &property.Description, &property.Rules)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve property", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(property)
}

func (PropertyHandler *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	propertyID := params["id"]
	PropertyHandler.LoadProperties()
	var property Entities.Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE Property SET Name = ?, Address = ?, CreateTime = ?, Type = ?, Photos = ?, Description = ?, Rules = ? WHERE PropertyID = ?`
	_, err = PropertyHandler.db.Exec(query, property.Name, property.Address, property.CreateTime, property.Type, property.Photos, property.Description, property.Rules, propertyID)
	if err != nil {
		http.Error(w, "Failed to update property", http.StatusInternalServerError)
		return
	}
	PropertyHandler.LoadProperties()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Property updated successfully")
}

func (PropertyHandler *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	propertyID := params["id"]
	PropertyHandler.LoadProperties()
	query := `DELETE FROM Property WHERE PropertyID = ?`
	_, err := PropertyHandler.db.Exec(query, propertyID)
	if err != nil {
		http.Error(w, "Failed to delete property", http.StatusInternalServerError)
		return
	}
	PropertyHandler.LoadProperties()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Property deleted successfully")
}
