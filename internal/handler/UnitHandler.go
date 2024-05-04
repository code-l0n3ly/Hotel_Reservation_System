package Handlers

import (
	"database/sql"
	"net/http"
	"strconv"
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
	rows, err := UnitHandler.db.Query(`SELECT UnitID, PropertyID, Name, RentalPrice, Description, Rating, OccupancyStatus, StructuralProperties, CreateTime FROM Unit`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var unit Entities.Unit
		if err := rows.Scan(&unit.UnitID, &unit.PropertyID, &unit.Name, &unit.RentalPrice, &unit.Description, &unit.Rating, &unit.OccupancyStatus, &unit.StructuralProperties, &createTime); err != nil {
			return err
		}
		unit.CreateTime, err = time.Parse("2006-01-02 15:04:05", string(createTime))
		if err != nil {
			return err
		}
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

	query := `INSERT INTO Unit (PropertyID, Name, RentalPrice, Description, Rating, OccupancyStatus, StructuralProperties) VALUES (?, ?, ?, ?, ?, ?, ?)`
	tx, err := UnitHandler.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create unit"})
		return
	}
	result, err := tx.Exec(query, unit.PropertyID, unit.Name, unit.RentalPrice, unit.Description, unit.Rating, unit.OccupancyStatus, unit.StructuralProperties)
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
	err = tx.Commit()
	id, _ := result.LastInsertId()
	unit.UnitID = strconv.FormatInt(id, 10)
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
	var unit Entities.Unit
	var createTime []byte

	query := `SELECT UnitID, PropertyID, Name, RentalPrice, Description, Rating, OccupancyStatus, StructuralProperties, CreateTime FROM Unit WHERE UnitID = ?`
	err := UnitHandler.db.QueryRow(query, unitID).Scan(&unit.UnitID, &unit.PropertyID, &unit.Name, &unit.RentalPrice, &unit.Description, &unit.Rating, &unit.OccupancyStatus, &unit.StructuralProperties, &createTime)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve unit" + err.Error()})
		return
	}

	unit.CreateTime, err = time.Parse("2006-01-02 15:04:05", string(createTime))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve unit" + err.Error()})
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
	OldInfoUnit := UnitHandler.cache[unitID]
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
	if NewInfoUnit.PropertyID != "" {
		OldInfoUnit.PropertyID = NewInfoUnit.PropertyID
	}
	if NewInfoUnit.Name != "" {
		OldInfoUnit.Name = NewInfoUnit.Name
	}
	if NewInfoUnit.Description != "" {
		OldInfoUnit.Description = NewInfoUnit.Description
	}
	if NewInfoUnit.OccupancyStatus != "" {
		OldInfoUnit.OccupancyStatus = NewInfoUnit.OccupancyStatus
	}
	if NewInfoUnit.StructuralProperties != "" {
		OldInfoUnit.StructuralProperties = NewInfoUnit.StructuralProperties
	}
	if NewInfoUnit.Rating != 0 {
		OldInfoUnit.Rating = NewInfoUnit.Rating
	}
	if NewInfoUnit.RentalPrice != 0 {
		OldInfoUnit.RentalPrice = NewInfoUnit.RentalPrice
	}
	query := `UPDATE Unit SET Name = ?, Description = ?, RentalPrice = ?, Rating = ?, OccupancyStatus = ?, StructuralProperties = ? WHERE UnitID = ?`
	_, err = UnitHandler.db.Exec(query, OldInfoUnit.Name, OldInfoUnit.Description, OldInfoUnit.RentalPrice, OldInfoUnit.Rating, OldInfoUnit.OccupancyStatus, OldInfoUnit.StructuralProperties, OldInfoUnit.UnitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update unit" + err.Error()})
		return
	}

	rows, err := UnitHandler.db.Query(`SELECT Image FROM Images WHERE UnitID = ?`, unitID)
	if err != nil {
	} else {
		var existingImages []string
		for rows.Next() {
			var image string
			if err := rows.Scan(&image); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update unit" + err.Error()})
				return
			}
			existingImages = append(existingImages, image)
		}

		if len(existingImages) > 0 {
			_, err = UnitHandler.db.Exec(`DELETE FROM Images WHERE UnitID = ?`, unitID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update unit" + err.Error()})
				return
			}
		}
	}
	if NewInfoUnit.Images != nil {
		for _, image := range NewInfoUnit.Images {
			_, err = UnitHandler.db.Exec(`INSERT INTO Images (UnitID, Image, Type) VALUES (?, ?, ?, ?, ?)`, OldInfoUnit.UnitID, image, "Unit")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update unit" + err.Error()})
				return
			}
		}
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
