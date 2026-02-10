package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func Certificate(r *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewCertificateHandler(cfg)

	r.GET("/verify", h.VerifyCertificate)

	adminRoutes := r.Group("")
	adminRoutes.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin"}))
	{
		adminRoutes.POST("/issue-all/:eventID", h.IssueAllCertificates)
		adminRoutes.POST("/issue/:registerID", h.IssueOneCertificate)
	}
	userRoutes := r.Group("")
	userRoutes.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"default"}))
	{
		userRoutes.GET("/:id/download", h.Download)
		userRoutes.GET("/me", h.ME)
	}

}
