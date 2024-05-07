package Handlers

import (
	"database/sql"
	"net/http"
	"time"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	db                *sql.DB
	ReviewIdReference int64
	cache             map[string]Entities.Review // Cache to hold users in memory
}

func NewReviewHandler(db *sql.DB) *ReviewHandler {
	return &ReviewHandler{
		db:    db,
		cache: make(map[string]Entities.Review),
	}
}

func (ReviewHandler *ReviewHandler) LoadReviews() error {
	rows, err := ReviewHandler.db.Query(`SELECT ReviewID, UserID, UnitID, Rating, Comment, CreateTime FROM Review`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime []byte
		var review Entities.Review
		if err := rows.Scan(&review.ReviewID, &review.UserID, &review.UnitID, &review.Rating, &review.Comment, &createTime); err != nil {
			return err
		}
		//fmt.Println(review)
		review.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		ReviewHandler.cache[review.ReviewID] = review
	}
	return rows.Err()
}

func (ReviewHandler *ReviewHandler) CreateReview(c *gin.Context) {
	var review Entities.Review
	ReviewHandler.LoadReviews()
	err := c.BindJSON(&review)
	if err != nil {

		c.JSON(http.StatusBadRequest, Response{
			Status:  "success",
			Message: "Review created successfully",
			Data:    review,
		})
		return
	}

	err = review.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	query := `INSERT INTO Review (UserID, UnitID, Rating, Comment) VALUES (?, ?, ?, ?)`
	_, err = ReviewHandler.db.Exec(query, review.UserID, review.UnitID, review.Rating, review.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to create review",
		})
		return
	}
	ReviewHandler.LoadReviews()
	c.JSON(http.StatusCreated, ReviewHandler.cache[review.ReviewID]) // Respond with the created review object
}

func (ReviewHandler *ReviewHandler) GetReview(c *gin.Context) {
	reviewID := c.Param("id")

	var review Entities.Review

	// Load the cache
	ReviewHandler.LoadReviews()

	// Check if the review is in the cache
	if cachedReview, ok := ReviewHandler.cache[reviewID]; ok {
		review = cachedReview
	} else {
		// If the review is not in the cache, return an error
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Review not found",
		})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (ReviewHandler *ReviewHandler) UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")
	ReviewHandler.LoadReviews()
	var review Entities.Review
	err := c.BindJSON(&review)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	err = review.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	query := `UPDATE Review SET `
	setValues := []interface{}{}
	if review.UserID != "" {
		query += `UserID = ?, `
		setValues = append(setValues, review.UserID)
	}
	if review.UnitID != "" {
		query += `UnitID = ?, `
		setValues = append(setValues, review.UnitID)
	}
	if review.Rating != 0 {
		query += `Rating = ?, `
		setValues = append(setValues, review.Rating)
	}
	if review.Comment != "" {
		query += `Comment = ?, `
		setValues = append(setValues, review.Comment)
	}

	// Remove the last comma and space
	query = query[:len(query)-2]

	query += ` WHERE ReviewID = ?`
	setValues = append(setValues, reviewID)

	_, err = ReviewHandler.db.Exec(query, setValues...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to update review",
		})
		return
	}
	ReviewHandler.LoadReviews()
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Review updated successfully",
	})
}
func (ReviewHandler *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	ReviewHandler.LoadReviews()
	query := `DELETE FROM Review WHERE ReviewID = ?`
	_, err := ReviewHandler.db.Exec(query, reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to delete review",
		})
		return
	}
	ReviewHandler.LoadReviews()
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Review deleted successfully",
	})
}
