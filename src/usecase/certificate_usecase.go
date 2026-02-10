package usecase

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/pdf"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
)

type CertificateUsecase struct {
	config            *config.Config
	certificateRepo   repository.CertificateRepository
	registrationsRepo repository.RegistrationRepository
	pdfClient         *pdf.Client
}

type CertTemplateData struct {
	Username      string
	EventName     string
	Duration      string
	Date          string
	CertificateID string
}

func NewCertificateUsecase(cfg *config.Config,
	certificateRepository repository.CertificateRepository,
	registrationRepository repository.RegistrationRepository,
	pdfClient *pdf.Client) *CertificateUsecase {
	return &CertificateUsecase{
		config:            cfg,
		certificateRepo:   certificateRepository,
		registrationsRepo: registrationRepository,
		pdfClient:         pdfClient}
}

func (uc *CertificateUsecase) IssueEventCertificate(ctx context.Context, eventId int) error {
	participants, err := uc.registrationsRepo.GetAllAttendedByEventId(ctx, eventId)
	if err != nil {
		return err
	}
	if len(participants) == 0 {
		return &service_errors.ServiceError{EndUserMessage: "No attended participants found"}
	}

	var certsToCreate []model.Certificate
	for _, p := range participants {
		certsToCreate = append(certsToCreate, model.Certificate{
			RegistrationId: p.Id,
			EventId:        p.EventId,
			TrackingCode:   uc.generateTrackingCode(),
			Status:         model.Pending,
		})
	}
	if len(certsToCreate) == 0 {
		return &service_errors.ServiceError{EndUserMessage: "No eligible participants for certificate"}
	}

	createdCertificates, err := uc.certificateRepo.BulkCreate(ctx, certsToCreate)
	if err != nil {
		log.Printf("[ERROR] Failed to bulk create certificates: %v", err)
		return err
	}

	go uc.issueCertificatesAsync(eventId, createdCertificates)

	return nil
}

func (uc *CertificateUsecase) IssueRegistrationCertificate(ctx context.Context, registrationId int) error {
	certificate, err := uc.certificateRepo.GetByRegistrationId(ctx, registrationId)
	if err != nil || certificate == nil {
		tCode := uc.generateTrackingCode()
		registration, err := uc.registrationsRepo.FindById(ctx, registrationId)
		if err != nil {
			return err
		}
		cert := model.Certificate{
			RegistrationId: registration.Id,
			EventId:        registration.EventId,
			TrackingCode:   tCode,
			Status:         model.Pending,
		}
		certificate, err = uc.certificateRepo.Create(ctx, cert)
		if err != nil {
			return err
		}
	}
	return uc.issueCertificate(ctx, *certificate)
}

func (uc *CertificateUsecase) issueCertificatesAsync(eventId int, certificates []model.Certificate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var failedCount int32

	for _, c := range certificates {
		sem <- struct{}{}
		wg.Add(1)

		go func(cert model.Certificate) {
			defer func() {
				<-sem
				wg.Done()
			}()

			if err := uc.issueCertificate(ctx, cert); err != nil {
				atomic.AddInt32(&failedCount, 1)
				log.Printf("[ERROR] Failed to issue certificate for event %d: %v", cert.Id, err)
			}
		}(c)
	}
	wg.Wait()
	log.Printf("[INFO] Successfully issued %d certificates", len(certificates))
}

func (uc *CertificateUsecase) issueCertificate(ctx context.Context, cert model.Certificate) error {
	certificate, err := uc.certificateRepo.GetById(ctx, cert.Id)
	if err != nil {
		return err
	}
	if certificate.Status != model.Pending {
		return &service_errors.ServiceError{EndUserMessage: "certificate already issued"}
	}
	//HTML
	html := uc.generateHTML(certificate.Registration, certificate.TrackingCode)
	if html == "" {
		log.Printf("[ERROR] HTML generation failed for certID: %d", certificate.Id)
		return &service_errors.ServiceError{EndUserMessage: "HTML generation failed"}
	}

	//PDF
	pdfBytes, err := uc.pdfClient.HTMLToPDF(html)
	if err != nil {
		log.Printf("[ERROR] Gotenberg failed for certID %d: %v", certificate.Id, err)
		return &service_errors.ServiceError{EndUserMessage: "Gotenberg failed"}
	}

	//Save
	filePath, err := uc.savePDF(pdfBytes, certificate.RegistrationId)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: "Save PDF failed"}
	}

	//Issue in DB
	file := &model.FileRef{
		Path: filePath,
		Size: int64(len(pdfBytes)),
		Mime: "application/pdf",
	}
	err = uc.certificateRepo.MarkAsIssued(ctx, certificate.Id, file)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: "Mark certificate as issued failed"}
	}
	return nil
}

func (uc *CertificateUsecase) generateHTML(participant model.Registration, trackingCode string) string {
	tmplPath := "./template/certificate.html"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		log.Println("Error reading template file:", err)
		return ""
	}

	data := CertTemplateData{
		Username:      participant.User.FirstName + " " + participant.User.LastName,
		EventName:     participant.Event.Title,
		Duration:      "12",
		Date:          common.ToShamsiString(participant.Event.Date),
		CertificateID: trackingCode,
	}

	t, err := template.New("certificate").Parse(string(tmplBytes))
	if err != nil {
		log.Println("Error parsing template:", err)
		return ""
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Println("Error executing template:", err)
		return ""
	}

	return buf.String()
}

func (uc *CertificateUsecase) savePDF(pdfBytes []byte, registrationID int) (string, error) {
	uploadDir := "storage/certificates"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("[ERROR] Failed to create directory %s: %v", uploadDir, err)
			return "", fmt.Errorf("could not create storage directory: %w", err)
		}
	}

	fileName := fmt.Sprintf("cert_%d_%d.pdf", registrationID, time.Now().Unix())
	filePath := filepath.Join(uploadDir, fileName)

	err := os.WriteFile(filePath, pdfBytes, 0644)
	if err != nil {
		log.Printf("[ERROR] Failed to write PDF file for RegID %d: %v", registrationID, err)
		return "", fmt.Errorf("failed to save file to disk: %w", err)
	}

	log.Printf("[SUCCESS] Certificate saved locally at: %s", filePath)

	return filePath, nil
}

func (uc *CertificateUsecase) generateTrackingCode() string {
	// return uuid.New().String()
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
