package handler

import (
	"errors"
	"net/http"
	"strconv"

	reqDto "github.com/farzadamr/event-manager-api/api/dto"
	"github.com/farzadamr/event-manager-api/api/helper"
	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/usecase"
	"github.com/farzadamr/event-manager-api/usecase/dto"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleUsecase *usecase.RoleUsecase
	config      *config.Config
}

func NewRoleHandler(cfg *config.Config) *RoleHandler {
	roleRepo := dependency.GetRoleRepository()
	return &RoleHandler{
		config:      cfg,
		roleUsecase: usecase.NewRoleUsecase(cfg, roleRepo),
	}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req reqDto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}

	roleDto, err := common.TypeConverter[dto.CreateRole](req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	err = h.roleUsecase.Create(c.Request.Context(), roleDto)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true))
}

func (h *RoleHandler) GetById(c *gin.Context) {
	id := c.Param("eventID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("event id required")))
		return
	}
	roleId, convErr := strconv.Atoi(id)
	if convErr != nil || roleId < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid event id")))
		return
	}

	role, err := h.roleUsecase.GetById(c.Request.Context(), roleId)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(role, true))
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var req reqDto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, err))
		return
	}

	err := h.roleUsecase.Update(c.Request.Context(), req.Id, req.DisplayName)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true))
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(roles, true))
}
