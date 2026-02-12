package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewUserHandler(cfg)

	router.POST("/login-by-student-number", h.LoginByStudentNumber)
	router.POST("/register-by-student-number", h.RegisterByStudentNumber)
	router.POST("/refresh-token", h.RefreshToken)
	adminRoutes := router.Group("")
	adminRoutes.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"admin"}))
	{
		adminRoutes.GET("/list", h.GetList)
		adminRoutes.GET("/:id", h.GetById)
		adminRoutes.POST(":id/assign-roles", h.AssignRoles)
		adminRoutes.POST(":id/revoke-roles", h.RevokeRoles)
	}
	userRoutes := router.Group("")
	userRoutes.Use(middleware.Authentication(cfg), middleware.Authorization([]string{"default"}))
	{
		userRoutes.PUT("/edit-profile", h.EditProfile)
	}

}
