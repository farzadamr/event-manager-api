package repository

import (
	"context"

	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"gorm.io/gorm"
)

type CertificateRepository struct {
	database *gorm.DB
	preloads []database.PreloadEntity
}

func NewCertificateRepository(preloads []database.PreloadEntity) *CertificateRepository {
	return &CertificateRepository{database: database.GetDb(), preloads: preloads}
}

func (cr *CertificateRepository) Create(ctx context.Context, r model.Certificate) (model.Certificate, error) {
	err := cr.database.WithContext(ctx).Create(&r).Error
	if err != nil {
		return model.Certificate{}, err
	}
	return r, nil
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
