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
	var createTime []byte
	rows, err := PropertyHandler.db.Query(`SELECT PropertyID, OwnerID, Name, Description, Address, CreateTime FROM Property`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var property Entities.Property
		if err := rows.Scan(&property.PropertyID, &property.OwnerID, &property.Name, &property.Description, &property.Address, createTime); err != nil {
			return err
		}
		property.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		PropertyHandler.cache[property.PropertyID] = property
	}
	PropertyHandler.SetHighestPropertyID()
	return rows.Err()
}

func (PropertyHandler *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	PropertyHandler.LoadProperties()
	var property Entities.Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	propertyID := PropertyHandler.GenerateUniquePropertyID()
	property.PropertyID = propertyID

	query := `INSERT INTO Property (PropertyID, OwnerID, Name, Address, CreateTime, Type, Description, Rules) VALUES (?, ?, ?, ?, NOW(), ?, ?, ?)`
	_, err = PropertyHandler.db.Exec(query, property.PropertyID, property.OwnerID, property.Name, property.Address, property.Type, property.Description, property.Rules)
	if err != nil {
		response := Response{
			Status:  "Failed to create property",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Insert images
	for _, image := range property.Photos {
		_, err = PropertyHandler.db.Exec(`INSERT INTO Images (PropertyID, Type, Image) VALUES (?, 'Proof', ?)`, propertyID, image)
		if err != nil {
			response := Response{
				Status:  "Failed to create property",
				Message: err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	PropertyHandler.LoadProperties()
	response := Response{
		Status:  "success",
		Message: "Property created successfully",
		Data:    PropertyHandler.cache[propertyID],
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (PropertyHandler *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	PropertyHandler.LoadProperties()
	var property Entities.Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	query := `UPDATE Property SET OwnerID = ?, Name = ?, Address = ?, Type = ?, Description = ?, Rules = ? WHERE PropertyID = ?`
	_, err = PropertyHandler.db.Exec(query, property.OwnerID, property.Name, property.Address, property.Type, property.Description, property.Rules, property.PropertyID)
	if err != nil {
		response := Response{
			Status:  "Failed to update property",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if there are existing images
	var count int
	err = PropertyHandler.db.QueryRow(`SELECT COUNT(*) FROM Images WHERE PropertyID = ? AND Type = 'Proof'`, property.PropertyID).Scan(&count)
	if err != nil {
		response := Response{
			Status:  "Failed to update property",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if count > 0 {
		// Delete existing images
		_, err = PropertyHandler.db.Exec(`DELETE FROM Images WHERE PropertyID = ? AND Type = 'Proof'`, property.PropertyID)
		if err != nil {
			response := Response{
				Status:  "Failed to update property",
				Message: err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Insert new images
	for _, image := range property.Photos {
		_, err = PropertyHandler.db.Exec(`INSERT INTO Images (PropertyID, Type, Image) VALUES (?, 'Proof', ?)`, property.PropertyID, image)
		if err != nil {
			response := Response{
				Status:  "Failed to update property",
				Message: err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response := Response{
		Status:  "success",
		Message: "Property updated successfully",
		Data:    property,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (PropertyHandler *PropertyHandler) GetProperties(w http.ResponseWriter, r *http.Request) {
	PropertyHandler.LoadProperties()
	var properties []Entities.Property
	for _, property := range PropertyHandler.cache {
		properties = append(properties, property)
	}

	response := Response{
		Status:  "success",
		Message: "Properties retrieved successfully",
		Data:    properties,
	}

	json.NewEncoder(w).Encode(response)
}

func (PropertyHandler *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	PropertyHandler.LoadProperties()
	params := mux.Vars(r)
	propertyID := params["id"]

	property, ok := PropertyHandler.cache[propertyID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Property retrieved successfully",
		Data:    property,
	}

	json.NewEncoder(w).Encode(response)
}

func (PropertyHandler *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	PropertyHandler.LoadProperties()
	params := mux.Vars(r)
	propertyID := params["id"]

	query := `DELETE FROM Property WHERE PropertyID = ?`
	_, err := PropertyHandler.db.Exec(query, propertyID)
	if err != nil {
		response := Response{
			Status:  "Failed to delete property",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete images
	_, err = PropertyHandler.db.Exec(`DELETE FROM Images WHERE PropertyID = ?`, propertyID)
	if err != nil {
		response := Response{
			Status:  "Failed to delete property",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	PropertyHandler.LoadProperties()
	response := Response{
		Status:  "success",
		Message: "Property deleted successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get all properties by UserID
func (PropertyHandler *PropertyHandler) GetPropertiesByUserID(w http.ResponseWriter, r *http.Request) {
	PropertyHandler.LoadProperties()
	params := mux.Vars(r)
	userID := params["id"]

	var properties []Entities.Property
	for _, property := range PropertyHandler.cache {
		if property.OwnerID == userID {
			properties = append(properties, property)
		}
	}

	response := Response{
		Status:  "success",
		Message: "Properties retrieved successfully",
		Data:    properties,
	}

	json.NewEncoder(w).Encode(response)
}
