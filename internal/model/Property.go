package model

import (
	"errors"
	"time"
)

// Property represents the 'Property' table in your database.
type Property struct {
	PropertyID  string    `json:"propertyID"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	CreateTime  time.Time `json:"createTime"`
	Type        string    `json:"type"`
	Photos      [][]byte  `json:"images,omitempty"`
	OwnerID     string    `json:"ownerID"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"` // Assuming JSON data as a string; adjust according to your needs
}

func (p *Property) Validate() error {
	if p.PropertyID == "" {
		return errors.New("PropertyID is required")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Address == "" {
		return errors.New("address is required")
	}
	if p.Type == "" {
		return errors.New("type is required")
	}

	if p.Description == "" {
		return errors.New("description is required")
	}
	if p.Rules == "" {
		return errors.New("rules are required")
	}
	return nil
}

func (p *Property) HasRules() bool {
	return p.Rules != ""
}
