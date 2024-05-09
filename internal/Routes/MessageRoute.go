package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(router *gin.Engine, MessageHandler *handler.MessageHandler) {
	router.POST("/message/send", MessageHandler.SendMessage)
	//router.GET("/chat/:id", MessageHandler.GetChat)
	router.GET("/chat/:id", MessageHandler.GetChatByID)
}
