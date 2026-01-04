package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/farzadamr/event-manager-api/api/helper"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/usecase"
	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	registerUsecase *usecase.RegisterEventUsecase
	config          *config.Config
}

func NewRegisterHandler(cfg *config.Config) *RegisterHandler {
	userRepository := dependency.GetUserRepository()
	registerUsecase := usecase.NewRegisterEventUsecase(cfg, dependency.GetRegisterEventRepository(), userRepository)
	return &RegisterHandler{registerUsecase: registerUsecase, config: cfg}
}

func (h *RegisterHandler) RegisterEvent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("event id required")))
		return
	}
	eventID, err := strconv.Atoi(id)
	if err != nil || eventID < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid event id")))
		return
	}

	userId := int(c.Value(constant.UserIdKey).(float64))
	err = h.registerUsecase.RegisterForEvent(c, eventID, userId)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err), helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true))
}

func (h *RegisterHandler) CancelRegisteration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("event id required")))
		return
	}
	eventID, err := strconv.Atoi(id)
	if err != nil || eventID < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid event id")))
		return
	}
	userId := int(c.Value(constant.UserIdKey).(float64))
	err = h.registerUsecase.CancelRegisteration(c, eventID, userId)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err), helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true))
}
