package handler

import (
	"net/http"
	"strconv"

	"github.com/farzadamr/event-manager-api/api/dto"
	"github.com/farzadamr/event-manager-api/api/helper"
	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/usecase"
	usecaseDto "github.com/farzadamr/event-manager-api/usecase/dto"
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
			helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}

	token, err := h.userUsecase.LoginByStudentNumber(c, req.StudentNumber, req.Password)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
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

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true))
}

func (h *UserHandler) RegisterByStudentNumber(c *gin.Context) {
	var req dto.RegisterUserByStudentNumberRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}
	err = h.userUsecase.RegisterByStudentNumber(c, req.ToRegisterUserByStudentNumber())
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true))

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
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(token, true))
}

func (h *UserHandler) GetList(c *gin.Context) {
	roleName := c.DefaultQuery("role", "default")
	if err := h.userUsecase.ValidateRoleName(roleName); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}

	pn := c.DefaultQuery("pageNumber", "1")
	pageNumber, err := strconv.Atoi(pn)
	if err != nil || pageNumber < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	ps := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(ps)
	if err != nil || pageSize < 1 || pageSize > 50 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	pagination := filter.PaginationInput{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	pagedResult, err := h.userUsecase.GetUserListByRoleName(c.Request.Context(), roleName, pagination)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(pagedResult, true))
}

func (h *UserHandler) EditProfile(c *gin.Context) {
	var req dto.UpdateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}
	updateUserDto, err := common.TypeConverter[usecaseDto.UpdateUserDto](req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	err = h.userUsecase.Update(c, updateUserDto)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true))
}
