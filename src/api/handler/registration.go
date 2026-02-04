package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/farzadamr/event-manager-api/api/helper"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/usecase"
	"github.com/farzadamr/event-manager-api/usecase/dto"
	"github.com/gin-gonic/gin"
)

type RegistrationsHandler struct {
	registrationUsecase *usecase.RegistrationUseCase
	config              *config.Config
}

func NewRegistrationHandler(cfg *config.Config) *RegistrationsHandler {
	userRepo := dependency.GetUserRepository()
	registerRepo := dependency.GetRegisterEventRepository()
	return &RegistrationsHandler{registrationUsecase: usecase.NewRegistrationUseCase(cfg, registerRepo, userRepo), config: cfg}
}

func (h *RegistrationsHandler) GetRegistrations(c *gin.Context) {
	id := c.Param("eventID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("event id required")))
		return
	}
	eventID, err := strconv.Atoi(id)
	if err != nil || eventID < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid event id")))
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

	pagedResult, err := h.registrationUsecase.GetRegistrations(c, eventID, pagination)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(pagedResult, true))
}

func (h *RegistrationsHandler) GetAttendanceList(c *gin.Context) {
	id := c.Param("eventID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("event id required")))
		return
	}
	eventID, err := strconv.Atoi(id)
	if err != nil || eventID < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid event id")))
		return
	}

	pn := c.DefaultQuery("pageNumber", "1")
	pageNumber, err := strconv.Atoi(pn)
	if err != nil || pageNumber < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	ps := c.DefaultQuery("pageSize", "50")
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

	pagedResult, err := h.registrationUsecase.AttendanceList(c, eventID, pagination)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(pagedResult, true))
}

func (h *RegistrationsHandler) UpdateAttendanceList(c *gin.Context) {
	var req []dto.AttendanceRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}

	if err = h.registrationUsecase.UpdateAttendanceList(c, req); err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true))
}
