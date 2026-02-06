package usecase

import (
	"context"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/pdf"
)

type CertificateUsecase struct {
	config          *config.Config
	certificateRepo repository.CertificateRepository
	pdfClient       *pdf.Client
}

func NewCertificateUsecase(cfg *config.Config, certificateRepository repository.CertificateRepository, pdfClient *pdf.Client) *CertificateUsecase {
	return &CertificateUsecase{config: cfg, certificateRepo: certificateRepository, pdfClient: pdfClient}
}

func (uc *CertificateUsecase) IssueEventCertificate(ctx context.Context, eventId int) error {

}
