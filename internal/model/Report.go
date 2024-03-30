package model

import (
	"database/sql"
	"errors"
	"time"
)

// Report represents the 'Report' table in your database.
type Report struct {
	ReportID   string         `json:"reportID"`
	UserID     string         `json:"userID"`
	Type       sql.NullString `json:"type,omitempty"`
	CreateTime time.Time      `json:"createTime,omitempty"`
	Data       string         `json:"data"` // Assuming JSON data as a string; adjust according to your needs
}

func (r *Report) Validate() error {
	if r.ReportID == "" {
		return errors.New("ReportID is required")
	}
	if r.UserID == "" {
		return errors.New("UserID is required")
	}
	if r.Data == "" {
		return errors.New("data is required")
	}
	return nil
}

func (r *Report) HasType() bool {
	return r.Type.Valid
}
