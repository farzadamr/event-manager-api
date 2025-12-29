package router

import (
	"github.com/farzadamr/event-manager-api/api/handler"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewUserHandler(cfg)

	router.POST("/login-by-student-number", h.LoginByStudentNumber)
	router.POST("/register-by-student-number", h.RegisterByStudentNumber)
	router.POST("/refresh-token", h.RefreshToken)
}
