package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	db                *sql.DB
	ReviewIdReference int64
	cache             map[string]Entities.Review // Cache to hold users in memory
}

func NewReviewHandler(db *sql.DB) *ReviewHandler {
	return &ReviewHandler{
		db:                db,
		ReviewIdReference: 0,
		cache:             make(map[string]Entities.Review),
	}
}

func (ReviewHandler *ReviewHandler) GenerateUniqueReviewID() string {
	ReviewHandler.ReviewIdReference++
	return fmt.Sprintf("%d", ReviewHandler.ReviewIdReference)
}

func (ReviewHandler *ReviewHandler) SetHighestReviewID() {
	highestID := int64(0)
	for _, review := range ReviewHandler.cache {
		reviewID, err := strconv.ParseInt(review.ReviewID, 10, 64)
		if err != nil {
			continue // Skip if the ReviewID is not a valid integer
		}
		if reviewID > highestID {
			highestID = reviewID
		}
	}
	ReviewHandler.ReviewIdReference = highestID
}

func (ReviewHandler *ReviewHandler) LoadReviews() error {
	rows, err := ReviewHandler.db.Query(`SELECT ReviewID, UserID, UnitID, Rating, Comment, CreateTime FROM Review`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var review Entities.Review
		if err := rows.Scan(&review.ReviewID, &review.UserID, &review.UnitID, &review.Rating, &review.Comment, &review.CreateTime); err != nil {
			return err
		}
		//fmt.Println(review)
		ReviewHandler.cache[review.ReviewID] = review
	}
	ReviewHandler.SetHighestReviewID()
	return rows.Err()
}

func (ReviewHandler *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var review Entities.Review
	ReviewHandler.LoadReviews()
	ReviewHandler.GenerateUniqueReviewID()
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = review.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO Review (ReviewID, UserID, UnitID, Rating, Comment, CreateTime) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = ReviewHandler.db.Exec(query, review.ReviewID, review.UserID, review.UnitID, review.Rating, review.Comment, review.CreateTime)
	if err != nil {
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}
	ReviewHandler.LoadReviews()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ReviewHandler.cache[review.ReviewID]) // Respond with the created review object
}

func (ReviewHandler *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	reviewID := params["id"]

	var review Entities.Review
	query := `SELECT ReviewID, UserID, UnitID, Rating, Comment, CreateTime FROM Review WHERE ReviewID = ?`
	err := ReviewHandler.db.QueryRow(query, reviewID).Scan(&review.ReviewID, &review.UserID, &review.UnitID, &review.Rating, &review.Comment, &review.CreateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve review", http.StatusInternalServerError)
		return
	}
	ReviewHandler.LoadReviews()
	json.NewEncoder(w).Encode(review)
}

func (ReviewHandler *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	reviewID := params["id"]
	ReviewHandler.LoadReviews()
	var review Entities.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = review.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE Review SET UserID = ?, UnitID = ?, Rating = ?, Comment = ?, CreateTime = ? WHERE ReviewID = ?`
	_, err = ReviewHandler.db.Exec(query, review.UserID, review.UnitID, review.Rating, review.Comment, review.CreateTime, reviewID)
	if err != nil {
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}
	ReviewHandler.LoadReviews()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Review updated successfully")
}

func (ReviewHandler *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	reviewID := params["id"]
	ReviewHandler.LoadReviews()
	query := `DELETE FROM Review WHERE ReviewID = ?`
	_, err := ReviewHandler.db.Exec(query, reviewID)
	if err != nil {
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}
	ReviewHandler.LoadReviews()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Review deleted successfully")
}
