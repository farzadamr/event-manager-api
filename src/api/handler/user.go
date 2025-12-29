package handler

import (
	"net/http"

	"github.com/farzadamr/event-manager-api/api/dto"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	config      *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	userUsecase := usecase.NewUserUsecase(cfg, dependency.GetUserRepository())
	return &UserHandler{
		userUsecase: userUsecase,
		config:      cfg,
	}
}

func (h *UserHandler) RegisterByStudentNumber(c *gin.Context) {
	var req dto.RegisterUserByStudentNumberRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}
	err = h.userUsecase.RegisterByStudentNumber(c, req.ToRegisterUserByStudentNumber())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}
