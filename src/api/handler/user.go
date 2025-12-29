package handler

import (
	"net/http"

	"github.com/farzadamr/event-manager-api/api/dto"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase  *usecase.UserUsecase
	tokenUsecase *usecase.TokenUsecase
	config       *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	userUsecase := usecase.NewUserUsecase(cfg, dependency.GetUserRepository())
	tokenUsecase := usecase.NewTokenUsecase(cfg)
	return &UserHandler{
		userUsecase:  userUsecase,
		tokenUsecase: tokenUsecase,
		config:       cfg,
	}
}
func (h *UserHandler) LoginByStudentNumber(c *gin.Context) {
	var req dto.LoginByStudentNumberRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"Message": "Validation Error"})
		return
	}

	token, err := h.userUsecase.LoginByStudentnumber(c, req.StudentNumber, req.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constant.RefreshTokenCookieName,
		Value:    token.RefreshToken,
		MaxAge:   int(h.config.JWT.RefreshTokenExpireDuration * 60),
		Path:     "/",
		Domain:   h.config.Server.Domain,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusCreated, token)
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

func (h *UserHandler) RefreshToken(c *gin.Context) {
	token, err := h.tokenUsecase.RefreshToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constant.RefreshTokenCookieName,
		Value:    token.RefreshToken,
		MaxAge:   int(h.config.JWT.RefreshTokenExpireDuration * 60),
		Path:     "/",
		Domain:   h.config.Server.Domain,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	c.JSON(http.StatusOK, token)
}
