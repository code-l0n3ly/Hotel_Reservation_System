package model

import (
	"errors"
	"time"
)

type Unit struct {
	UnitID               string    `json:"unitID"`
	AddressID            string    `json:"addressID"`
	Name                 string    `json:"name,omitempty"` // Optional field
	Images               [][]byte  `json:"images,omitempty"`
	Description          string    `json:"description,omitempty"`
	Rating               float32   `json:"rating,omitempty"`
	PropertyID           string    `json:"propertyID"`
	RentalPrice          int       `json:"rentalPrice"`
	OccupancyStatus      string    `json:"occupancyStatus"`
	StructuralProperties string    `json:"structuralProperties"` // Assuming JSON data as a string; adjust according to your needs
	CreateTime           time.Time `json:"createTime"`
	Address              Address   `json:"address"`
}

func (u *Unit) Validate() error {
	if u.PropertyID == "" {
		return errors.New("PropertyID is required")
	}
	if u.RentalPrice < 0 {
		return errors.New("RentalPrice must be greater than 0")
	}
	if u.OccupancyStatus == "" {
		return errors.New("OccupancyStatus is required Expected values: 'occupied' or 'available'")
	}
	if u.StructuralProperties == "" {
		return errors.New("StructuralProperties is required")
	}
	return nil
}

func (u *Unit) IsOccupied() bool {
	return u.OccupancyStatus == "occupied"
}

func (u *Unit) IsAvailable() bool {
	return u.OccupancyStatus == "available"
}
