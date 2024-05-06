package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(router *gin.Engine, MessageHandler *handler.MessageHandler) {
	router.POST("/message/create", MessageHandler.CreateMessage)
	router.GET("/message/:id", MessageHandler.GetMessage)
	router.PUT("/message/:id", MessageHandler.UpdateMessage)
	router.DELETE("/message/:id", MessageHandler.DeleteMessage)
}
