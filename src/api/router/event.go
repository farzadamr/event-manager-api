package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func Event(router *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewEventHandler(cfg)
	rh := handler.NewRegisterHandler(cfg)
	router.GET("/:id", h.GetEventById)
	router.GET("", h.GetEvents)

	adminRoutes := router.Group("")
	adminRoutes.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin"}))
	{
		adminRoutes.POST("", h.Create)
		adminRoutes.PATCH("/:id", h.Update)
		adminRoutes.PATCH("/:id/status", h.ChangeEventStatus)
		adminRoutes.DELETE("/:id", h.Delete)
	}

	userProtected := router.Group("")
	userProtected.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin", "teacher", "default"}))
	{
		userProtected.POST("/:id/register", rh.RegisterEvent)
		userProtected.DELETE("/:id/register", rh.CancelRegisteration)
	}
}
