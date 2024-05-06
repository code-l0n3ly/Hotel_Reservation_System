package Handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gin-gonic/gin"
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

func (UnitHandler *UnitHandler) LoadUnits() error {
	query := `
        SELECT 
            u.UnitID, u.PropertyID, u.AddressID, u.Name, u.RentalPrice, u.Description, u.Rating, u.OccupancyStatus, u.StructuralProperties, u.CreateTime,
            a.AddressID, a.Country, a.City, a.State, a.Street, a.PostalCode, a.AdditionalNumber, a.MapLocation, a.Latitude, a.Longitude
        FROM 
            Unit u
        LEFT JOIN 
            Address a ON u.AddressID = a.AddressID
    `
	rows, err := UnitHandler.db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var unit Entities.Unit
		var address Entities.Address
		if err := rows.Scan(&unit.UnitID, &unit.PropertyID, &unit.AddressID, &unit.Name, &unit.RentalPrice, &unit.Description, &unit.Rating, &unit.OccupancyStatus, &unit.StructuralProperties, &createTime, &address.AddressID, &address.Country, &address.City, &address.State, &address.Street, &address.PostalCode, &address.AdditionalNumber, &address.MapLocation, &address.Latitude, &address.Longitude); err != nil {
			fmt.Println(err.Error())
		}
		unit.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))

		// Save the grouped object to the cache
		unit.Address = address
		fmt.Println(unit)
		UnitHandler.cache[unit.UnitID] = unit
	}
	return rows.Err()
}

func (UnitHandler *UnitHandler) CreateUnit(c *gin.Context) {
	var unit Entities.Unit
	UnitHandler.LoadUnits()
	err := c.BindJSON(&unit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	err = unit.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	tx, _ := UnitHandler.db.Begin()
	addressQuery := `INSERT INTO Address (Country, City, State, Street, PostalCode, AdditionalNumber, MapLocation, Latitude, Longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	addressResult, err := tx.Exec(addressQuery, unit.Address.Country, unit.Address.City, unit.Address.State, unit.Address.Street, unit.Address.PostalCode, unit.Address.AdditionalNumber, unit.Address.MapLocation, unit.Address.Latitude, unit.Address.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create address" + err.Error()})
		return
	}
	AddressID, _ := addressResult.LastInsertId()
	unit.AddressID = strconv.FormatInt(AddressID, 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create unit"})
		return
	}
	query := `INSERT INTO Unit (PropertyID, AddressID, Name, RentalPrice, Description, Rating, OccupancyStatus, StructuralProperties) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, unit.PropertyID, unit.AddressID, unit.Name, unit.RentalPrice, unit.Description, unit.Rating, unit.OccupancyStatus, unit.StructuralProperties)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create unit" + err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve unit ID" + err.Error()})
		return
	}

	if unit.Images != nil {
		for _, image := range unit.Images {
			_, err = tx.Exec(`INSERT INTO Images (UnitID, Image, Type) VALUES (?, ?, ?)`, unit.UnitID, image, "Unit")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to insert image" + err.Error()})
				return
			}
		}
	}
	id, _ := result.LastInsertId()
	unit.UnitID = strconv.FormatInt(id, 10)
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create unit" + err.Error()})
		return
	}
	UnitHandler.LoadUnits()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Unit created successfully", "data": UnitHandler.cache[unit.UnitID]})
}

func (UnitHandler *UnitHandler) GetUnit(c *gin.Context) {
	unitID := c.Param("id")
	UnitHandler.LoadUnits()

	unit, ok := UnitHandler.cache[unitID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Units retrieved successfully", "data": unit})
}

func (UnitHandler *UnitHandler) GetUnits(c *gin.Context) {
	UnitHandler.LoadUnits()
	var units []Entities.Unit
	for _, unit := range UnitHandler.cache {
		units = append(units, unit)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Units retrieved successfully", "data": units})
}

func (UnitHandler *UnitHandler) UpdateUnit(c *gin.Context) {
	unitID := c.Param("id")
	UnitHandler.LoadUnits()
	var NewInfoUnit Entities.Unit
	OldInfoUnit, ok := UnitHandler.cache[unitID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Unit not found"})
		return
	}

	err := c.BindJSON(&NewInfoUnit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	errValid := NewInfoUnit.Validate()
	if errValid != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errValid})
		return
	}

	// Prepare dynamic SQL for unit update
	updateUnitQuery := "UPDATE Unit SET "
	updateUnitParams := []interface{}{}
	fields := []string{}

	if NewInfoUnit.PropertyID != "" {
		fields = append(fields, "PropertyID = ?")
		updateUnitParams = append(updateUnitParams, NewInfoUnit.PropertyID)
		OldInfoUnit.PropertyID = NewInfoUnit.PropertyID
	}
	if NewInfoUnit.Name != "" {
		fields = append(fields, "Name = ?")
		updateUnitParams = append(updateUnitParams, NewInfoUnit.Name)
		OldInfoUnit.Name = NewInfoUnit.Name
	}
	if NewInfoUnit.Description != "" {
		fields = append(fields, "Description = ?")
		updateUnitParams = append(updateUnitParams, NewInfoUnit.Description)
		OldInfoUnit.Description = NewInfoUnit.Description
	}
	if NewInfoUnit.OccupancyStatus != "" {
		fields = append(fields, "OccupancyStatus = ?")
		updateUnitParams = append(updateUnitParams, NewInfoUnit.OccupancyStatus)
		OldInfoUnit.OccupancyStatus = NewInfoUnit.OccupancyStatus
	}
	if NewInfoUnit.StructuralProperties != "" {
		fields = append(fields, "StructuralProperties = ?")
		updateUnitParams = append(updateUnitParams, NewInfoUnit.StructuralProperties)
		OldInfoUnit.StructuralProperties = NewInfoUnit.StructuralProperties
	}
	if NewInfoUnit.Rating != 0 {
		fields = append(fields, "Rating = ?")
		updateUnitParams = append(updateUnitParams, NewInfoUnit.Rating)
		OldInfoUnit.Rating = NewInfoUnit.Rating
	}
	if NewInfoUnit.RentalPrice != 0 {
		fields = append(fields, "RentalPrice = ?")
		updateUnitParams = append(updateUnitParams, NewInfoUnit.RentalPrice)
		OldInfoUnit.RentalPrice = NewInfoUnit.RentalPrice
	}

	updateUnitQuery += strings.Join(fields, ", ") + " WHERE UnitID = ?"
	updateUnitParams = append(updateUnitParams, OldInfoUnit.UnitID)

	_, err = UnitHandler.db.Exec(updateUnitQuery, updateUnitParams...)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update unit" + err.Error()})
		return
	}

	// Handle images
	if NewInfoUnit.Images != nil {
		_, err = UnitHandler.db.Exec(`DELETE FROM Images WHERE UnitID = ?`, unitID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update unit" + err.Error()})
			return
		}

		for _, image := range NewInfoUnit.Images {
			_, err = UnitHandler.db.Exec(`INSERT INTO Images (UnitID, Image, Type) VALUES (?, ?, ?)`, OldInfoUnit.UnitID, image, "Unit")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update unit" + err.Error()})
				return
			}
		}
	}

	// Update the address related to the unit
	addressQuery := "UPDATE Address SET "
	params := []interface{}{}

	if NewInfoUnit.Address.Country != "" {
		addressQuery += "Country = ?, "
		params = append(params, NewInfoUnit.Address.Country)
	}
	if NewInfoUnit.Address.City != "" {
		addressQuery += "City = ?, "
		params = append(params, NewInfoUnit.Address.City)
	}
	if NewInfoUnit.Address.State != "" {
		addressQuery += "State = ?, "
		params = append(params, NewInfoUnit.Address.State)
	}
	if NewInfoUnit.Address.Street != "" {
		addressQuery += "Street = ?, "
		params = append(params, NewInfoUnit.Address.Street)
	}
	if NewInfoUnit.Address.PostalCode != "" {
		addressQuery += "PostalCode = ?, "
		params = append(params, NewInfoUnit.Address.PostalCode)
	}
	if NewInfoUnit.Address.AdditionalNumber != "" {
		addressQuery += "AdditionalNumber = ?, "
		params = append(params, NewInfoUnit.Address.AdditionalNumber)
	}
	if NewInfoUnit.Address.MapLocation != "" {
		addressQuery += "MapLocation = ?, "
		params = append(params, NewInfoUnit.Address.MapLocation)
	}
	if NewInfoUnit.Address.Latitude != "" {
		addressQuery += "Latitude = ?, "
		params = append(params, NewInfoUnit.Address.Latitude)
	}
	if NewInfoUnit.Address.Longitude != "" {
		addressQuery += "Longitude = ?, "
		params = append(params, NewInfoUnit.Address.Longitude)
	}

	// Remove the trailing comma and space
	addressQuery = addressQuery[:len(addressQuery)-2]

	addressQuery += " WHERE AddressID = ?"
	params = append(params, NewInfoUnit.AddressID)

	_, err = UnitHandler.db.Exec(addressQuery, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Unit updated successfully", "Data": OldInfoUnit})
}

func (UnitHandler *UnitHandler) DeleteUnit(c *gin.Context) {
	unitID := c.Param("id")
	UnitHandler.LoadUnits()
	query := `DELETE FROM Unit WHERE UnitID = ?`
	_, err := UnitHandler.db.Exec(query, unitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to delete unit" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Unit deleted successfully", "data": UnitHandler.cache[unitID]})
	UnitHandler.LoadUnits()
}

// GetAllUnits : Gets all the units that are available
func (UnitHandler *UnitHandler) GetAllAvailableUnits(c *gin.Context) {
	UnitHandler.LoadUnits()
	var units []Entities.Unit
	for _, unit := range UnitHandler.cache {
		if unit.OccupancyStatus == "Available" {
			units = append(units, unit)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Units retrieved successfully", "data": units})
}

func (UnitHandler *UnitHandler) GetAllOccupiedUnits(c *gin.Context) {
	UnitHandler.LoadUnits()
	var units []Entities.Unit
	for _, unit := range UnitHandler.cache {
		if unit.OccupancyStatus == "Occupied" {
			units = append(units, unit)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Units retrieved successfully", "data": units})
}

// function that search units by name, let it search if there is a unit with the exact name and then units that contain the name
func (UnitHandler *UnitHandler) SearchUnitsByName(c *gin.Context) {
	UnitHandler.LoadUnits()
	var unit Entities.Unit
	err := c.BindJSON(&unit)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	var units []Entities.Unit
	for _, unitx := range UnitHandler.cache {
		if unitx.Name == unit.Name {
			units = append(units, unitx)
		}
	}
	for _, unitx := range UnitHandler.cache {
		if strings.Contains(unit.Name, unit.Name) {
			units = append(units, unitx)
		}
	}
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Units retrieved successfully",
		Data:    units,
	})
}

// function that search units by address, let it search if there is a unit with the exact address and then units that contain the address
func (UnitHandler *UnitHandler) SearchUnitsByAddress(c *gin.Context) {
	UnitHandler.LoadUnits()
	var Address Entities.Address
	err := c.BindJSON(&Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	var units []Entities.Unit
	for _, unit := range UnitHandler.cache {
		if unit.Address.Street == Address.Street || unit.Address.City == Address.City || unit.Address.State == Address.State || unit.Address.Country == Address.Country || unit.Address.PostalCode == Address.PostalCode {
			units = append(units, unit)
		}
	}
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Units retrieved successfully",
		Data:    units,
	})
}
