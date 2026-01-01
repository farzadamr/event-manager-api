package repository

import (
	"context"
	"errors"

	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"gorm.io/gorm"
)

type RegistrationRepository struct {
	database *gorm.DB
	preloads []database.PreloadEntity
}

func NewRegistrationRepository(preloads []database.PreloadEntity) *RegistrationRepository {
	return &RegistrationRepository{database: database.GetDb(), preloads: preloads}
}

func (r *RegistrationRepository) Create(ctx context.Context, re model.Registration) (model.Registration, error) {
	err := r.database.WithContext(ctx).Create(&re).Error
	if err != nil {
		return model.Registration{}, err
	}
	return re, nil
}

func (r *RegistrationRepository) FindByEventIDAndUserID(ctx context.Context, eventID, userID int) (model.Registration, error) {
	var rg model.Registration
	q := "user_id = ? and event_id = ? " + softDeleteExp
	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)
	err := db.
		Where(q, userID, eventID).
		First(&rg).
		Error
	if err != nil {
		return model.Registration{}, err
	}
	return rg, nil
}

func (r *RegistrationRepository) ListByEventID(ctx context.Context, eventID int, pagination filter.PaginationInput) (int64, []model.Registration, error) {
	q := "event_id = ?" + softDeleteExp
	var totalRows int64

	if err := r.database.WithContext(ctx).
		Model(&model.Registration{}).
		Where(q, eventID).
		Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	var items []model.Registration
	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)

	offset := pagination.GetOffset()
	limit := pagination.GetPageSize()

	if err := db.
		Where(q, eventID).
		Offset(offset).
		Limit(limit).
		Find(&items).Error; err != nil {
		return 0, nil, err
	}
	return totalRows, items, nil
}

func (r *RegistrationRepository) ListByUserID(ctx context.Context, userId int, pagination filter.PaginationInput) (int64, []model.Registration, error) {
	q := "user_id = ?" + softDeleteExp
	var totalRows int64

	if err := r.database.WithContext(ctx).
		Model(&model.Registration{}).
		Where(q, userId).
		Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	var items []model.Registration
	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)

	offset := pagination.GetOffset()
	limit := pagination.GetPageSize()

	if err := db.
		Where(q, userId).
		Offset(offset).
		Limit(limit).
		Find(&items).Error; err != nil {
		return 0, nil, err
	}
	return totalRows, items, nil
}

func (r *RegistrationRepository) CancelByUser(ctx context.Context, eventID, userID int) error {
	q := "event_id = ? and user_id = ? and status = ? " + softDeleteExp
	result := r.database.WithContext(ctx).
		Model(&model.Registration{}).
		Where(q, eventID, userID, model.StatusRegistered).
		Update("status", model.StatusCancelledByUser)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("registration not found or already cancelled")
	}

	return nil
}

func (r *RegistrationRepository) CancelByEvent(ctx context.Context, eventID int) error {
	q := "event_id = ? and status = ? " + softDeleteExp
	result := r.database.WithContext(ctx).
		Model(&model.Registration{}).
		Where(q, eventID, model.StatusRegistered).
		Update("status", model.StatusCancelledByEvent)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (r *RegistrationRepository) UpdateAttendanceStatus(ctx context.Context, registrationId int, status model.AttendanceStatus) error {
	if valid := status == model.Present ||
		status == model.Absent ||
		status == model.NotCheckedIn; !valid {
		return errors.New("invalid attendance status")
	}

	q := "id = ? " + softDeleteExp

	result := r.database.WithContext(ctx).
		Model(&model.Registration{}).
		Where(q, registrationId).
		Update("attendance_status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
