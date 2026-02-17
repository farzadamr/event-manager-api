package repository

import (
	"context"
	"time"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"gorm.io/gorm"
)

type CertificateRepository struct {
	database *gorm.DB
	logger   logging.Logger
	preloads []database.PreloadEntity
}

func NewCertificateRepository(cfg *config.Config, preloads []database.PreloadEntity) *CertificateRepository {
	return &CertificateRepository{
		database: database.GetDb(),
		logger:   logging.NewLogger(cfg),
		preloads: preloads,
	}
}

func (cr *CertificateRepository) Create(ctx context.Context, r model.Certificate) (*model.Certificate, error) {

	db := cr.database.WithContext(ctx)

	if err := db.Create(&r).Error; err != nil {
		cr.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}

	var result model.Certificate
	preloadDB := database.Preload(db, cr.preloads)

	if err := preloadDB.
		First(&result, r.Id).Error; err != nil {

		cr.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}

	return &result, nil
}

func (cr *CertificateRepository) BulkCreate(ctx context.Context, certs []model.Certificate) ([]model.Certificate, error) {

	db := cr.database.WithContext(ctx)

	// insert
	if err := db.Create(&certs).Error; err != nil {
		cr.logger.Error(logging.Postgres, logging.BulkInsert, err.Error(), nil)
		return nil, err
	}

	// collect ids (FIXED)
	ids := make([]int, 0, len(certs))
	for _, c := range certs {
		ids = append(ids, c.Id)
	}

	var result []model.Certificate

	preloadDB := database.Preload(db, cr.preloads)

	if err := preloadDB.
		Where("id IN ?", ids).
		Find(&result).Error; err != nil {

		cr.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}

	return result, nil
}

func (cr *CertificateRepository) MarkAsIssued(ctx context.Context, id int, file *model.FileRef, metadata *model.CertificateMetadata) error {

	tx := cr.database.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var certificate model.Certificate
	if err := tx.
		First(&certificate, id).Error; err != nil {

		tx.Rollback()
		return &service_errors.ServiceError{
			EndUserMessage: service_errors.RecordNotFound,
		}
	}

	if certificate.RegistrationId == 0 {
		tx.Rollback()
		return &service_errors.ServiceError{
			EndUserMessage: service_errors.UnExpectedError,
		}
	}

	now := time.Now()
	update := map[string]interface{}{
		"issued_at": now,
		"status":    model.Issued,
	}
	if file != nil {
		update["file_path"] = file.Path
		update["file_size"] = file.Size
		update["file_mime"] = file.Mime
	}
	if metadata != nil {
		update["metadata_user_name"] = metadata.UserName
		update["metadata_english_name"] = metadata.EnglishName
		update["metadata_date"] = metadata.Date
		update["metadata_duration"] = metadata.Duration
		update["metadata_event_name"] = metadata.EventName
	}

	if err := tx.Model(&model.Certificate{}).
		Where("id = ?", id).
		Updates(update).Error; err != nil {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}

	if err := tx.Model(&model.Registration{}).
		Where("id = ?", certificate.RegistrationId).
		Update("status", model.StatusIssueCertificate).Error; err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (cr *CertificateRepository) GetById(ctx context.Context, id int) (model.Certificate, error) {
	certificate := new(model.Certificate)

	db := cr.database.WithContext(ctx)
	db = database.Preload(db, cr.preloads)

	err := db.Where("id = ?", id).
		First(&certificate).
		Error
	if err != nil {
		return *certificate, err
	}
	return *certificate, nil
}

func (cr *CertificateRepository) GetByFilter(ctx context.Context, eventId int, req filter.PaginationInput) (int64, []model.Certificate, error) {
	var items []model.Certificate
	var totalRows int64 = 0
	db := cr.database.WithContext(ctx)
	db = database.Preload(db, cr.preloads)
	query := db.
		Model(&model.Certificate{}).
		Joins("JOIN registrations r ON certificates.registration_id = r.id").
		Where("r.event_id = ?", eventId)
	if err := query.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	offset := req.GetOffset()
	limit := req.GetPageSize()

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return 0, nil, err
	}
	return totalRows, items, nil
}

func (cr *CertificateRepository) GetAllByEventId(ctx context.Context, eventId int) ([]model.Certificate, error) {
	var certificates []model.Certificate

	db := cr.database.WithContext(ctx)
	db = database.Preload(db, cr.preloads)

	err := db.Where("event_id = ?", eventId).Find(&certificates).Error
	if err != nil {
		return nil, err
	}
	return certificates, nil
}

func (cr *CertificateRepository) GetByRegistrationId(ctx context.Context, registrationId int) (*model.Certificate, error) {
	certificate := new(model.Certificate)
	db := cr.database.WithContext(ctx)
	db = database.Preload(db, cr.preloads)
	err := db.Where("registration_id = ?", registrationId).First(&certificate).Error
	if err != nil {
		return certificate, err
	}
	return certificate, nil
}

func (cr *CertificateRepository) GetByUserIdByFilter(ctx context.Context, userId int, req filter.PaginationInput) (int64, []model.Certificate, error) {
	var items []model.Certificate
	var totalRows int64 = 0

	query := cr.database.WithContext(ctx).
		Model(&model.Certificate{}).
		Joins("JOIN registrations r ON certificates.registration_id = r.id").
		Where("r.user_id = ?", userId).
		Where("certificates.status = ?", model.Issued)

	err := query.Count(&totalRows).Error
	if err != nil {
		return 0, nil, err
	}

	offset := req.GetOffset()
	limit := req.GetPageSize()
	err = query.
		Preload("Registration.Event").
		Offset(offset).
		Limit(limit).
		Order("certificates.created_at DESC").
		Find(&items).Error
	if err != nil {
		return 0, nil, err
	}
	return totalRows, items, nil
}

func (cr *CertificateRepository) VerifyCertificate(ctx context.Context, trackingCode string) (model.Certificate, error) {
	var certificate model.Certificate
	db := cr.database.WithContext(ctx)
	db = database.Preload(db, cr.preloads)

	err := db.
		Where("tracking_code = ? AND status = ?", trackingCode, model.Issued).
		First(&certificate).Error
	if err != nil {
		cr.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return certificate, err
	}
	return certificate, nil

}
