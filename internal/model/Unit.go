package model

import (
	"errors"
	"time"
)

type Unit struct {
	UnitID               string    `json:"unitID"`
	OwnerName            string    `json:"ownerName,omitempty"` // Optional field
	AddressID            string    `json:"addressID"`
	Name                 string    `json:"name,omitempty"` // Optional field
	Images               [][]byte  `json:"images,omitempty"`
	Description          string    `json:"description,omitempty"`
	Rating               float32   `json:"rating,omitempty"`
	PropertyID           string    `json:"propertyID"`
	RentalPrice          int       `json:"rentalPrice"`
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

	if u.StructuralProperties == "" {
		return errors.New("StructuralProperties is required")
	}
	return nil
}
