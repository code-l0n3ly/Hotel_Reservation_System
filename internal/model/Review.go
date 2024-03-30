package model

import (
	"database/sql"
	"errors"
	"time"
)

// Review represents the 'Review' table in your database.
type Review struct {
	ReviewID   string         `json:"reviewID"`
	UserID     string         `json:"userID"`
	UnitID     string         `json:"unitID"`
	Rating     int            `json:"rating"`
	Comment    sql.NullString `json:"comment,omitempty"`
	CreateTime time.Time      `json:"createTime"`
}

func (r *Review) Validate() error {
	if r.ReviewID == "" {
		return errors.New("ReviewID is required")
	}
	if r.UserID == "" {
		return errors.New("UserID is required")
	}
	if r.UnitID == "" {
		return errors.New("UnitID is required")
	}
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	return nil
}

func (r *Review) HasComment() bool {
	return r.Comment.Valid
}
