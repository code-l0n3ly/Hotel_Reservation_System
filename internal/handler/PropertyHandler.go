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

type PropertyHandler struct {
	db    *sql.DB
	cache map[string]Entities.Property // Cache to hold properties in memory
}

func NewPropertyHandler(db *sql.DB) *PropertyHandler {
	return &PropertyHandler{
		db:    db,
		cache: make(map[string]Entities.Property),
	}
}

func (PropertyHandler *PropertyHandler) LoadProperties() error {
	query := `
        SELECT 
            p.PropertyID, p.OwnerID, p.AddressID,  p.Name, p.Description, p.Type, p.Rules, p.CreateTime,
            a.AddressID, a.Country, a.City, a.State, a.Street, a.PostalCode, a.AdditionalNumber, a.MapLocation, a.Latitude, a.Longitude
        FROM 
            Property p
        LEFT JOIN 
            Address a ON p.AddressID = a.AddressID
    `
	rows, err := PropertyHandler.db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var property Entities.Property
		var address Entities.Address
		if err := rows.Scan(&property.PropertyID, &property.OwnerID, &property.AddressID, &property.Name, &property.Description, &property.Type, &property.Rules, &createTime, &address.AddressID, &address.Country, &address.City, &address.State, &address.Street, &address.PostalCode, &address.AdditionalNumber, &address.MapLocation, &address.Latitude, &address.Longitude); err != nil {
			fmt.Println(err.Error())
		}
		property.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))

		// Save the grouped object to the cache
		property.Address = address
		fmt.Println(property)
		PropertyHandler.cache[property.PropertyID] = property
	}
	return rows.Err()
}

func (PropertyHandler *PropertyHandler) CreateProperty(c *gin.Context) {
	var property Entities.Property
	PropertyHandler.LoadProperties()

	err := c.BindJSON(&property)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message here": err.Error()})
		return
	}

	err = property.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	tx, err := PropertyHandler.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create property"})
		return
	}

	// Insert into Address table
	addressQuery := `INSERT INTO Address (Country, City, State, Street, PostalCode, AdditionalNumber, MapLocation, Latitude, Longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	addressResult, err := tx.Exec(addressQuery, property.Address.Country, property.Address.City, property.Address.State, property.Address.Street, property.Address.PostalCode, property.Address.AdditionalNumber, property.Address.MapLocation, property.Address.Latitude, property.Address.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create address" + err.Error()})
		return
	}
	AddressID, _ := addressResult.LastInsertId()
	property.AddressID = strconv.FormatInt(AddressID, 10)
	// Insert into Property table
	propertyQuery := `INSERT INTO Property (OwnerID, AddressID, Name, Description, Type, Rules) VALUES (?, ?, ?, ?, ?, ?)`
	propertyResult, err := tx.Exec(propertyQuery, property.OwnerID, property.AddressID, property.Name, property.Description, property.Type, property.Rules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create property" + err.Error()})
		return
	}
	propertyID, _ := propertyResult.LastInsertId()
	property.PropertyID = strconv.FormatInt(propertyID, 10)
	tx.Commit()
	PropertyHandler.LoadProperties()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Property created successfully", "data": PropertyHandler.cache[property.PropertyID]})
}

func (PropertyHandler *PropertyHandler) UpdateOrInsertProof(c *gin.Context) {
	// Get the PropertyID from the URL parameters
	PropertyID := c.Param("propertyID")

	// Get the Proof from the request body
	Proof, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get request body"})
		return
	}

	// Check if an image for this property and type is proof
	query := `SELECT COUNT(*) FROM Images WHERE PropertyID = ? AND Type = 'proof'`
	row := PropertyHandler.db.QueryRow(query, PropertyID)
	var count int
	err = row.Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		// If not, insert a new row
		insertQuery := `INSERT INTO Images (PropertyID, Type, Image) VALUES (?, 'proof', ?)`
		_, err := PropertyHandler.db.Exec(insertQuery, PropertyID, Proof)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// If yes, update the existing row
		updateQuery := `UPDATE Images SET Image = ? WHERE PropertyID = ? AND Type = 'proof'`
		_, err := PropertyHandler.db.Exec(updateQuery, Proof, PropertyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (PropertyHandler *PropertyHandler) GetProof(c *gin.Context) {
	// Get the PropertyID from the URL parameters
	PropertyID := c.Param("propertyID")

	// Execute the SQL query
	query := `SELECT Image FROM Images WHERE PropertyID = ? AND Type = 'proof'`
	row := PropertyHandler.db.QueryRow(query, PropertyID)
	var proof []byte
	err := row.Scan(&proof)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the proof as a response
	c.Data(http.StatusOK, "application/octet-stream", proof)
}

func (PropertyHandler *PropertyHandler) GetProperty(c *gin.Context) {
	ownerID := c.Param("id")
	PropertyHandler.LoadProperties()

	property, exists := PropertyHandler.cache[ownerID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Property retrieved successfully", "data": property})
}

func (PropertyHandler *PropertyHandler) GetProperties(c *gin.Context) {
	PropertyHandler.LoadProperties()
	var properties []Entities.Property
	for _, property := range PropertyHandler.cache {
		properties = append(properties, property)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Properties retrieved successfully", "data": properties})
}

func (PropertyHandler *PropertyHandler) UpdateProperty(c *gin.Context) {
	ownerID := c.Param("id")
	PropertyHandler.LoadProperties()

	var newInfoProperty Entities.Property
	oldInfoProperty := PropertyHandler.cache[ownerID]

	err := c.BindJSON(&newInfoProperty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	errValid := newInfoProperty.Validate()
	if errValid != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errValid})
		return
	}

	if newInfoProperty.Name != "" {
		oldInfoProperty.Name = newInfoProperty.Name
	}
	if newInfoProperty.Description != "" {
		oldInfoProperty.Description = newInfoProperty.Description
	}
	if newInfoProperty.Type != "" {
		oldInfoProperty.Type = newInfoProperty.Type
	}
	if newInfoProperty.Rules != "" {
		oldInfoProperty.Rules = newInfoProperty.Rules
	}

	query := `UPDATE Property SET Name = ?, Address = ?, Description = ?, Type = ?, Rules = ? WHERE OwnerID = ?`
	_, err = PropertyHandler.db.Exec(query, oldInfoProperty.Name, oldInfoProperty.Address, oldInfoProperty.Description, oldInfoProperty.Type, oldInfoProperty.Rules, oldInfoProperty.OwnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update property" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Property updated successfully", "Data": oldInfoProperty})
}

func (PropertyHandler *PropertyHandler) DeleteProperty(c *gin.Context) {
	ownerID := c.Param("id")
	PropertyHandler.LoadProperties()

	query := `DELETE FROM Property WHERE OwnerID = ?`
	_, err := PropertyHandler.db.Exec(query, ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to delete property" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Property deleted successfully", "data": PropertyHandler.cache[ownerID]})
	PropertyHandler.LoadProperties()
}

func (PropertyHandler *PropertyHandler) GetPropertiesByUserID(c *gin.Context) {
	userID := c.Param("id")
	PropertyHandler.LoadProperties()
	var properties []Entities.Property
	for _, property := range PropertyHandler.cache {
		if property.OwnerID == userID {
			properties = append(properties, property)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Properties retrieved successfully", "data": properties})
}

func (PropertyHandler *PropertyHandler) GetPropertiesByType(c *gin.Context) {
	propertyType := c.Param("type")
	PropertyHandler.LoadProperties()
	var properties []Entities.Property
	for _, property := range PropertyHandler.cache {
		if property.Type == propertyType {
			properties = append(properties, property)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Properties retrieved successfully", "data": properties})
}

func (PropertyHandler *PropertyHandler) GetUnitsByPropertyID(c *gin.Context) {
	propertyID := c.Param("id")

	query := `
    SELECT 
        u.UnitID, 
        u.PropertyID, 
        u.Name, 
        u.RentalPrice, 
        u.Description, 
        u.Rating, 
        u.OccupancyStatus, 
        u.StructuralProperties, 
        u.CreateTime,
        a.AddressID, 
        a.Country, 
        a.City, 
        a.State, 
        a.Street, 
        a.PostalCode, 
        a.AdditionalNumber, 
        a.MapLocation, 
        a.Latitude, 
        a.Longitude
    FROM 
        Unit u
    LEFT JOIN 
        Address a ON u.AddressID = a.AddressID
    WHERE 
        u.PropertyID = ?
`
	rows, err := PropertyHandler.db.Query(query, propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve units: " + err.Error()})
		return
	}
	defer rows.Close()

	var units []Entities.Unit
	for rows.Next() {
		var createTime []byte
		var unit Entities.Unit
		var address Entities.Address
		if err := rows.Scan(&unit.UnitID, &unit.PropertyID, &unit.Name, &unit.RentalPrice, &unit.Description, &unit.Rating, &unit.OccupancyStatus, &unit.StructuralProperties, &createTime, &address.AddressID, &address.Country, &address.City, &address.State, &address.Street, &address.PostalCode, &address.AdditionalNumber, &address.MapLocation, &address.Latitude, &address.Longitude); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to scan units: " + err.Error()})
			return
		}
		unit.CreateTime, err = time.Parse("2006-01-02 15:04:05", string(createTime))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to parse time: " + err.Error()})
			return
		}
		unit.Address = address
		units = append(units, unit)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Units retrieved successfully", "data": units})
}
