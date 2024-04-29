package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up the routes for the application
func RegisterUserRoutes(router *gin.Engine, UserHandler *handler.UserHandler) {
	router.POST("/users/create", UserHandler.CreateUserHandler)
	router.POST("/users/login", UserHandler.LoginHandler)
	router.GET("/users", UserHandler.GetUsersHandler)
	router.GET("/users/:id", UserHandler.GetUserHandler)
	router.PUT("/users/:id", UserHandler.UpdateUserHandler)
	router.DELETE("/users/:id", UserHandler.DeleteUserHandler)
}
