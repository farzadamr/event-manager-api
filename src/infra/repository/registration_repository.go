package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"gorm.io/gorm"
)

type RegistrationRepository struct {
	database *gorm.DB
	logger   logging.Logger
	preloads []database.PreloadEntity
}

func NewRegistrationRepository(cfg *config.Config, preloads []database.PreloadEntity) *RegistrationRepository {
	return &RegistrationRepository{
		database: database.GetDb(),
		logger:   logging.NewLogger(cfg),
		preloads: preloads}
}

func (r *RegistrationRepository) Create(ctx context.Context, re model.Registration) error {
	err := r.database.WithContext(ctx).Create(&re).Error
	if err != nil {
		r.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return err
	}
	return nil
}

func (r *RegistrationRepository) FindById(ctx context.Context, id int) (model.Registration, error) {
	var re model.Registration
	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)
	err := db.First(&re, id).Error
	if err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
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
	q := "event_id = ? " + softDeleteExp
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
		r.logger.Error(logging.Postgres, logging.Update, result.Error.Error(), nil)
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
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return err
	}

	return nil
}

func (r *RegistrationRepository) UpdateAttendanceList(ctx context.Context, attendanceList []model.AttendanceList) error {
	if len(attendanceList) == 0 {
		return errors.New("attendance list is empty")
	}

	valuesClause := make([]string, 0, len(attendanceList))
	args := make([]interface{}, 0, len(attendanceList)*2)

	for i, v := range attendanceList {
		// Cast id to BIGINT in SQL to avoid type ambiguity
		valuesClause = append(valuesClause, fmt.Sprintf("($%d::BIGINT, $%d)", i*2+1, i*2+2))
		args = append(args, v.Id, v.AttendanceStatus)
	}

	query := fmt.Sprintf(`
		UPDATE registrations
		SET attendance_status = v.status
		FROM (VALUES %s) AS v(id, status)
		WHERE registrations.id = v.id AND registrations.deleted_at IS NULL`,
		strings.Join(valuesClause, ","))

	if err := r.database.WithContext(ctx).Exec(query, args...).Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return err
	}
	return nil
}

func (r *RegistrationRepository) GetAllAttendedByEventId(ctx context.Context, eventID int) ([]model.Registration, error) {
	var registrations []model.Registration

	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)

	err := db.
		Where("event_id = ? AND attendance_status = ?", eventID, model.Present).
		Find(&registrations).
		Error
	if err != nil {
		return nil, err
	}
	return registrations, nil
}
