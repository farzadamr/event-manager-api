package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func Registration(router *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewRegistrationHandler(cfg)

	protected := router.Group("")
	protected.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin", "teacher"}))
	{
		protected.GET("/:eventID", h.GetRegistrations)
		protected.GET("/:eventID/attendance", h.GetAttendanceList)
		protected.PUT("/attendance", h.UpdateAttendanceList)
	}
}
