package routes

import (
	"github.com/fesbarbosa/melivendas-api/internal/adapters/input/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterItemRoutes(router *gin.Engine, itemHandler *handlers.ItemHandler) {
	v1 := router.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", itemHandler.Create)
			items.GET("", itemHandler.List)
			items.GET("/:id", itemHandler.GetByID)
			items.PUT("/:id", itemHandler.Update)
			items.DELETE("/:id", itemHandler.Delete)
		}
	}
}
