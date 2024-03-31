package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterReviewRoutes(router *mux.Router, ReviewHandler *handler.ReviewHandler) {
	router.HandleFunc("/reviews", ReviewHandler.CreateReview).Methods("POST")
	router.HandleFunc("/reviews/{id}", ReviewHandler.GetReview).Methods("GET")
	router.HandleFunc("/reviews/{id}", ReviewHandler.UpdateReview).Methods("PUT")
	router.HandleFunc("/reviews/{id}", ReviewHandler.DeleteReview).Methods("DELETE")
}
