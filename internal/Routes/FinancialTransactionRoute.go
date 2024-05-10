package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterFinancialTransactionRoutes(router *gin.Engine, FinancialTransactionHandler *handler.FinancialTransactionHandler) {
	router.POST("/financialTransaction/create", FinancialTransactionHandler.CreateTransaction)
	router.GET("/financialTransaction/:id", FinancialTransactionHandler.GetTransaction)
	router.PUT("/financialTransaction/:id", FinancialTransactionHandler.UpdateTransaction)
	router.DELETE("/financialTransaction/:id", FinancialTransactionHandler.DeleteTransaction)
}
