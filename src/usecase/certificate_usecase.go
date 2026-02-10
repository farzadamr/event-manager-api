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
	"time"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/pdf"
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

	go func() {
		sem := make(chan struct{}, 5)

		for _, p := range participants {
			sem <- struct{}{}

			go func(participant model.Registration) {
				defer func() { <-sem }()
				uc.issueOne(ctx, participant)
			}(p)
		}
	}()

	return nil
}

func (uc *CertificateUsecase) IssueRegistrationCertificate(ctx context.Context, registrationId int) error {
	registration, err := uc.registrationsRepo.FindById(ctx, registrationId)
	if err != nil {
		return err
	}

	go uc.issueOne(ctx, registration)
	return nil
}

func (uc *CertificateUsecase) issueOne(ctx context.Context, participant model.Registration) {
	tCode := uc.generateTrackingCode()

	//HTML
	html := uc.generateHTML(participant, tCode)
	if html == "" {
		log.Printf("[ERROR] HTML generation failed for User: %d", participant.UserId)
		return
	}

	//PDF
	pdfBytes, err := uc.pdfClient.HTMLToPDF(html)
	if err != nil {
		log.Printf("[ERROR] Gotenberg failed for RegID %d: %v", participant.Id, err)
		return
	}

	//Save
	filePath, err := uc.savePDF(pdfBytes, participant.Id)
	if err != nil {
		return
	}

	//DB
	err = uc.markAsIssued(ctx, participant.Id, filePath, tCode, pdfBytes)
	if err != nil {
		log.Printf("[ERROR] DB update failed for RegID %d: %v", participant.Id, err)
		return
	}

	log.Printf("[MONITOR] Certificate issued for %s (ID: %d)",
		participant.User.FirstName,
		participant.Id,
	)
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

	filePath := fmt.Sprintf("./storage/html/cert_%d_%d.html", participant.Id, time.Now().Unix())
	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		log.Println("Error writing HTML file:", err)
	} else {
		log.Println("HTML saved to:", filePath)
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

func (uc *CertificateUsecase) markAsIssued(ctx context.Context, registrationId int, filePath, trackingCode string, pdfBytes []byte) error {
	now := time.Now()
	cert := model.Certificate{
		RegistrationId: registrationId,
		IssuedAt:       &now,
		TrackingCode:   trackingCode,
		Pdf: model.FileRef{
			Path: filePath,
			Size: int64(len(pdfBytes)),
			Mime: "application/pdf",
		},
	}
	_, err := uc.certificateRepo.Create(ctx, cert)
	if err != nil {
		log.Printf("[ERROR] Failed to save certificate to DB: %v", err)
		return err
	}
	return nil
}
