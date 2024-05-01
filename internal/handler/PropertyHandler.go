package Handlers

import (
	"database/sql"
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
	rows, err := PropertyHandler.db.Query(`SELECT PropertyID, OwnerID, Name, Address, Description, Type, Rules, CreateTime FROM Property`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var property Entities.Property
		if err := rows.Scan(&property.PropertyID, &property.OwnerID, &property.Name, &property.Address, &property.Description, &property.Type, &property.Rules, &createTime); err != nil {
			return err
		}
		property.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		PropertyHandler.cache[property.PropertyID] = property
	}
	return rows.Err()
}

func (PropertyHandler *PropertyHandler) CreateProperty(c *gin.Context) {
	var property Entities.Property
	PropertyHandler.LoadProperties()

	err := c.BindJSON(&property)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	err = property.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := `INSERT INTO Property (OwnerID, Name, Address, Description, Type, Rules) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := PropertyHandler.db.Exec(query, property.OwnerID, property.Name, property.Address, property.Description, property.Type, property.Rules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create property" + err.Error()})
		return
	}
	id, _ := result.LastInsertId()
	property.PropertyID = strconv.FormatInt(id, 10)
	PropertyHandler.LoadProperties()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Property created successfully", "data": PropertyHandler.cache[property.OwnerID]})
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
	if newInfoProperty.Address != "" {
		oldInfoProperty.Address = newInfoProperty.Address
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

	query := `SELECT UnitID, PropertyID, Name, RentalPrice, Description, Rating, OccupancyStatus, StructuralProperties, CreateTime FROM Unit WHERE PropertyID = ?`
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
		if err := rows.Scan(&unit.UnitID, &unit.PropertyID, &unit.Name, &unit.RentalPrice, &unit.Description, &unit.Rating, &unit.OccupancyStatus, &unit.StructuralProperties, &createTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to scan units: " + err.Error()})
			return
		}
		unit.CreateTime, err = time.Parse("2006-01-02 15:04:05", string(createTime))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to parse time: " + err.Error()})
			return
		}
		units = append(units, unit)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Units retrieved successfully", "data": units})
}
