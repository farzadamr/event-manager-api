package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/farzadamr/event-manager-api/api/helper"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/dependency"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/pkg/pdf"
	"github.com/farzadamr/event-manager-api/usecase"
	"github.com/gin-gonic/gin"
)

type CertificateHandler struct {
	certificateUsecase *usecase.CertificateUsecase
	config             *config.Config
}

func NewCertificateHandler(cfg *config.Config) *CertificateHandler {
	gotenbergClient := pdf.NewGotenbergClient(cfg.Pdf.Url)
	registrationRepo := dependency.GetRegisterEventRepository()
	userRepo := dependency.GetUserRepository()
	certificateUsecase := usecase.NewCertificateUsecase(cfg, dependency.GetCertificateRepository(), registrationRepo, userRepo, gotenbergClient)
	return &CertificateHandler{config: cfg, certificateUsecase: certificateUsecase}
}

func (h *CertificateHandler) IssueAllCertificates(c *gin.Context) {
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

	if err = h.certificateUsecase.IssueEventCertificate(c, eventID); err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	result := map[string]string{
		"message": "Certificate issuing started",
	}
	c.JSON(http.StatusAccepted, helper.GenerateBaseResponse(result, true))
}

func (h *CertificateHandler) IssueOneCertificate(c *gin.Context) {
	id := c.Param("registerID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("register id required")))
		return
	}
	registerID, err := strconv.Atoi(id)
	if err != nil || registerID < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid register id")))
		return
	}

	if err = h.certificateUsecase.IssueRegistrationCertificate(c, registerID); err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusAccepted, helper.GenerateBaseResponse(nil, true))
}

func (h *CertificateHandler) Download(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("id required")))
		return
	}
	CertID, err := strconv.Atoi(id)
	if err != nil || CertID < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid id")))
		return
	}

	filePath, err := h.certificateUsecase.GetCertificateFile(c, CertID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	downloadName := fmt.Sprintf("certificate-%d.pdf", CertID)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadName))
	c.Header("Content-Type", "application/pdf")

	c.File(filePath)
}

func (h *CertificateHandler) ME(c *gin.Context) {
	pn := c.DefaultQuery("pageNumber", "1")
	pageNumber, err := strconv.Atoi(pn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	ps := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(ps)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	if pageNumber < 1 || pageSize < 1 || pageSize > 50 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, errors.New("pagination parameters invalid")))
		return
	}
	pagination := filter.PaginationInput{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	pagedResult, err := h.certificateUsecase.GetUserCertificates(c, pagination)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(pagedResult, true))
}

func (h *CertificateHandler) VerifyCertificate(c *gin.Context) {
	tCode := c.Query("ID")
	if strings.TrimSpace(tCode) == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("ID required")))
		return
	}
	res, _, err := h.certificateUsecase.VerifyCertificate(c.Request.Context(), tCode)
	if !res || err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	result := map[string]string{
		"message": "Certificate verified",
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(result, true))
}

func (h *CertificateHandler) GetByEventId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("id required")))
		return
	}
	eventId, err := strconv.Atoi(id)
	if err != nil || eventId < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, errors.New("invalid id")))
		return
	}

	pn := c.DefaultQuery("pageNumber", "1")
	pageNumber, err := strconv.Atoi(pn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	ps := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(ps)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}

	if pageNumber < 1 || pageSize < 1 || pageSize > 50 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, errors.New("pagination parameters invalid")))
		return
	}
	pagination := filter.PaginationInput{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	pagedResult, err := h.certificateUsecase.GetListByEventId(c, eventId, pagination)
	if err != nil {
		sc := helper.TranslateErrorToStatusCode(err)
		c.AbortWithStatusJSON(sc, helper.GenerateBaseResponseWithError(nil, false, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(pagedResult, true))
}
