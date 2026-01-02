package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/farzadamr/event-manager-api/api/dto"
	"github.com/farzadamr/event-manager-api/api/helper"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/usecase"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventUsecase *usecase.EventUsecase
	config       *config.Config
}

func NewEventHandler(cfg *config.Config) *EventHandler {
	userRepository := dependency.GetUserRepository()
	eventUsecase := usecase.NewEventUsecase(cfg, dependency.GetEventRepository(), userRepository)
	return &EventHandler{eventUsecase: eventUsecase, config: cfg}
}

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.CreateEventRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}

	err = h.eventUsecase.PublishEvent(c, req.ToCreateEvent())
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true))
}

func (h *EventHandler) Update(c *gin.Context) {
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
	var req dto.UpdateEventRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}
	req.Id = eventID
	err = h.eventUsecase.Update(c, req.ToUpdateEvent())
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true))
}

func (h *EventHandler) ChangeEventStatus(c *gin.Context) {
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

	if err = h.eventUsecase.ChangeEventStatus(c, eventID); err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true))
}

func (h *EventHandler) GetEvents(c *gin.Context) {
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
	pagedResult, err := h.eventUsecase.GetByFilter(c, pagination)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(pagedResult, true))
}

func (h *EventHandler) GetEventById(c *gin.Context) {
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

	event, err := h.eventUsecase.GetById(c, eventID)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(event, true))
}

func (h *EventHandler) Delete(c *gin.Context) {
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

	err = h.eventUsecase.Delete(c, eventID)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusNoContent, helper.GenerateBaseResponse(nil, true))
}
