package repository

import (
	"context"
	"time"

	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"gorm.io/gorm"
)

type CertificateRepository struct {
	database *gorm.DB
	preloads []database.PreloadEntity
}

func NewCertificateRepository(preloads []database.PreloadEntity) *CertificateRepository {
	return &CertificateRepository{database: database.GetDb(), preloads: preloads}
}

func (cr *CertificateRepository) Create(ctx context.Context, r model.Certificate) (*model.Certificate, error) {
	err := cr.database.WithContext(ctx).Create(&r).Error
	if err != nil {
		return &model.Certificate{}, err
	}
	return &r, nil
}

func (cr *CertificateRepository) BulkCreate(ctx context.Context, certs []model.Certificate) ([]model.Certificate, error) {
	err := cr.database.WithContext(ctx).Create(&certs).Error
	if err != nil {
		return nil, err
	}
	return certs, nil
}

func (cr *CertificateRepository) MarkAsIssued(ctx context.Context, id int, file *model.FileRef, metadata *model.CertificateMetadata) error {
	certificate := new(model.Certificate)
	if err := cr.database.WithContext(ctx).First(certificate, id).Error; err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	now := time.Now()
	certificate.IssuedAt = &now
	certificate.Pdf = file
	certificate.Metadata = metadata
	certificate.Status = model.Issued

	if err := cr.database.WithContext(ctx).Save(certificate).Error; err != nil {
		return err
	}
	return nil
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

	query := cr.database.WithContext(ctx).
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
		return certificate, err
	}
	return certificate, nil

}
