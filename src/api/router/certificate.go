package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func Certificate(r *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewCertificateHandler(cfg)

	adminRoutes := r.Group("")
	adminRoutes.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin"}))
	{
		adminRoutes.POST("/issue/:eventID", h.IssueCertificates)
	}
}
