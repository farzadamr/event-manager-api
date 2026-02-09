package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/farzadamr/event-manager-api/api/helper"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/dependency"
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
	certificateUsecase := usecase.NewCertificateUsecase(cfg, dependency.GetCertificateRepository(), registrationRepo, gotenbergClient)
	return &CertificateHandler{config: cfg, certificateUsecase: certificateUsecase}
}

func (h *CertificateHandler) IssueCertificates(c *gin.Context) {
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
