package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPropertyRoutes(router *gin.Engine, PropertyHandler *handler.PropertyHandler) {
	propertyRoutes := router.Group("/property")
	{
		propertyRoutes.POST("/create", PropertyHandler.CreateProperty)
		propertyRoutes.GET("/:id", PropertyHandler.GetProperty)
		propertyRoutes.GET("/", PropertyHandler.GetProperties)
		propertyRoutes.GET("/owner/:id", PropertyHandler.GetPropertiesByUserID)
		propertyRoutes.GET("/AllUnits/:id", PropertyHandler.GetUnitsByPropertyID)
		propertyRoutes.GET("/ByType/:type", PropertyHandler.GetPropertiesByType)
		propertyRoutes.PUT("/:id", PropertyHandler.UpdateProperty)
		propertyRoutes.DELETE("/:id", PropertyHandler.DeleteProperty)
		propertyRoutes.POST("/proof/add/:id", PropertyHandler.UpdateOrInsertProof)
		propertyRoutes.GET("/proof/get/:id", PropertyHandler.GetProof)
	}
}
