package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up the routes for the application
func RegisterUserRoutes(router *gin.Engine, UserHandler *handler.UserHandler) {
	users := router.Group("/users")
	{
		users.POST("/create", UserHandler.CreateUserHandler)
		users.POST("/login", UserHandler.LoginHandler)
		users.GET("/", UserHandler.GetUsersHandler)
		users.GET("/:id", UserHandler.GetUserHandler)
		users.PUT("/:id", UserHandler.UpdateUserHandler)
		users.DELETE("/:id", UserHandler.DeleteUserHandler)
	}
}
