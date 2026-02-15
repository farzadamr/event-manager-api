package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func Role(r *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewRoleHandler(cfg)

	routes := r.Group("")
	routes.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin"}))
	{
		routes.POST("/", h.CreateRole)
		routes.GET("/:id", h.GetById)
		routes.GET("/", h.GetAll)
		routes.PUT("/", h.UpdateRole)
	}
}
