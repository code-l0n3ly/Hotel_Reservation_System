package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(router *gin.Engine, ReviewHandler *handler.ReviewHandler) {
	router.POST("/reviews/create", ReviewHandler.CreateReview)
	router.GET("/reviews/:id", ReviewHandler.GetReview)
	router.PUT("/reviews/:id", ReviewHandler.UpdateReview)
	router.DELETE("/reviews/:id", ReviewHandler.DeleteReview)
	router.GET("/reviews/ByUnit/:id", ReviewHandler.GetReviewsByUnitID)
}
