package usecase

import (
	"bytes"
	"context"
	"errors"
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
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/farzadamr/event-manager-api/pkg/pdf"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"github.com/farzadamr/event-manager-api/usecase/dto"
	"gorm.io/gorm"
)

type CertificateUsecase struct {
	config            *config.Config
	logger            logging.Logger
	certificateRepo   repository.CertificateRepository
	registrationsRepo repository.RegistrationRepository
	userRepo          repository.UserRepository
	pdfClient         *pdf.Client
}

type CertTemplateData struct {
	Username      string
	EventName     string
	Duration      string
	StartDate     string
	CertificateID string
}

func NewCertificateUsecase(cfg *config.Config,
	certificateRepository repository.CertificateRepository,
	registrationRepository repository.RegistrationRepository,
	userRepo repository.UserRepository,
	pdfClient *pdf.Client) *CertificateUsecase {
	return &CertificateUsecase{
		config:            cfg,
		logger:            logging.NewLogger(cfg),
		certificateRepo:   certificateRepository,
		registrationsRepo: registrationRepository,
		userRepo:          userRepo,
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

	var ids []int
	for _, p := range participants {
		ids = append(ids, p.Id)
	}
	existingCerts, err := uc.certificateRepo.GetByRegistrationIds(ctx, ids)
	if err != nil {
		return err
	}

	existingMap := make(map[int]bool)
	for _, c := range existingCerts {
		existingMap[c.RegistrationId] = true
	}

	var certsToCreate []model.Certificate
	for _, p := range participants {
		if !existingMap[p.Id] {
			certsToCreate = append(certsToCreate, model.Certificate{
				RegistrationId: p.Id,
				EventId:        p.EventId,
				TrackingCode:   uc.generateTrackingCode(),
				Status:         model.Pending,
			})
		}
	}
	if len(certsToCreate) == 0 {
		if len(existingMap) != 0 {
			go uc.issueCertificatesAsync(existingCerts)
			return nil
		}
		return &service_errors.ServiceError{EndUserMessage: "No eligible participants for certificate"}
	}

	createdCertificates, err := uc.certificateRepo.BulkCreate(ctx, certsToCreate)
	if err != nil {
		return err
	}

	go uc.issueCertificatesAsync(createdCertificates)

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

func (uc *CertificateUsecase) GetListByEventId(ctx context.Context, eventId int, req filter.PaginationInput) (*filter.PagedList[dto.Certificate], error) {
	count, certificates, err := uc.certificateRepo.GetByFilter(ctx, eventId, req)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	dtoCertificates := dto.ToCertificateList(certificates)
	return filter.NewPagedList(&dtoCertificates, count, req.PageNumber, int64(req.PageSize)), nil
}

func (uc *CertificateUsecase) GetCertificateFile(ctx context.Context, certID int) (string, error) {
	userId, _ := ctx.Value(constant.UserIdKey).(float64)
	if userId == 0 {
		return "", &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	user, err := uc.userRepo.FetchUserInfoById(ctx, int(userId))
	if err != nil {
		return "", &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	cert, err := uc.certificateRepo.GetById(ctx, certID)
	if err != nil {
		return "", &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}

	if cert.Registration.UserId != user.Id {
		return "", &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	if cert.Status != model.Issued {
		return "", &service_errors.ServiceError{EndUserMessage: "certificate is not issued"}
	}

	return cert.Pdf.Path, nil
}

func (uc *CertificateUsecase) GetUserCertificates(ctx context.Context, req filter.PaginationInput) (*filter.PagedList[dto.UserCertificate], error) {
	userId := int(ctx.Value(constant.UserIdKey).(float64))
	if userId == 0 {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	user, err := uc.userRepo.FetchUserInfoById(ctx, userId)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	count, certificates, err := uc.certificateRepo.GetByUserIdByFilter(ctx, user.Id, req)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	dtoUserCertificates := dto.ToUserCertificateList(certificates)

	return filter.NewPagedList(&dtoUserCertificates, count, req.PageNumber, int64(req.PageSize)), nil

}

func (uc *CertificateUsecase) VerifyCertificate(ctx context.Context, trackingCode string) (bool, *model.Certificate, error) {
	result, err := uc.certificateRepo.VerifyCertificate(ctx, trackingCode)
	if err != nil {
		log.Printf("[ERROR] Failed to verify certificate: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		}
		return false, nil, &service_errors.ServiceError{EndUserMessage: service_errors.CertificateInvalid}
	}
	return true, &result, nil
}

func (uc *CertificateUsecase) issueCertificatesAsync(certificates []model.Certificate) {
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
				uc.logger.Warnf("[CERT] Failed to issue certificate for event %d: %v", cert.Id, err)
			}
		}(c)
	}
	wg.Wait()
	uc.logger.Infof("[CERT] Successfully issued %d certificates", len(certificates))
}

func (uc *CertificateUsecase) issueCertificate(ctx context.Context, certificate model.Certificate) error {
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
	duration := common.CalculateDurationString(certificate.Registration.Event.StartDate, certificate.Registration.Event.EndDate)
	metadata := &model.CertificateMetadata{
		UserName:    certificate.Registration.User.FirstName + " " + certificate.Registration.User.LastName,
		EventName:   certificate.Registration.Event.Title,
		EnglishName: certificate.Registration.User.EnglishName,
		Date:        common.ToShamsiString(certificate.Registration.Event.StartDate),
		Duration:    duration,
	}
	err = uc.certificateRepo.MarkAsIssued(ctx, certificate.Id, file, metadata)
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
	duration := common.CalculateDurationString(participant.Event.StartDate, participant.Event.EndDate)
	data := CertTemplateData{
		Username:      participant.User.FirstName + " " + participant.User.LastName,
		EventName:     participant.Event.Title,
		Duration:      duration,
		StartDate:     common.ToShamsiString(participant.Event.StartDate),
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
