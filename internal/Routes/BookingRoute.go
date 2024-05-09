package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterBookingRoutes(router *gin.Engine, BookingHandler *handler.BookingHandler) {
	bookingGroup := router.Group("/booking")
	{
		bookingGroup.POST("/create", BookingHandler.CreateBooking)
		bookingGroup.GET("/:id", BookingHandler.GetBooking)
		bookingGroup.PUT("/:id", BookingHandler.UpdateBooking)
		bookingGroup.DELETE("/:id", BookingHandler.DeleteBooking)
		bookingGroup.GET("/unit/:id", BookingHandler.GetActiveBookings)
		bookingGroup.GET("/user/:id", BookingHandler.GetBookingsByUserID)
	}
}
