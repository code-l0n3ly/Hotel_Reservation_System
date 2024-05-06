package Routes

import (
	handler "GraduationProject.com/m/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterUnitRoutes(router *gin.Engine, UnitHandler *handler.UnitHandler) {
	units := router.Group("/units")
	{
		units.POST("/create", UnitHandler.CreateUnit)
		units.GET("/:id", UnitHandler.GetUnit)
		units.GET("/", UnitHandler.GetUnits)
		units.PUT("/:id", UnitHandler.UpdateUnit)
		units.DELETE("/:id", UnitHandler.DeleteUnit)
		units.GET("/Available", UnitHandler.GetAllAvailableUnits)
		units.GET("/Occupied", UnitHandler.GetAllOccupiedUnits)
		units.POST("/SearchByName", UnitHandler.SearchUnitsByName)
		units.POST("/SearchByAddress", UnitHandler.SearchUnitsByAddress)
	}
}
