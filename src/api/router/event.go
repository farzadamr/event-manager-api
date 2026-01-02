package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func Event(router *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewEventHandler(cfg)

	router.GET("/:id", h.GetEventById)
	router.GET("", h.GetEvents)

	protected := router.Group("")
	protected.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin"}))
	{
		protected.POST("", h.Create)
		protected.PATCH("/:id", h.Update)
		protected.PATCH("/:id/status", h.ChangeEventStatus)
		protected.DELETE("/:id", h.Delete)
	}
}
