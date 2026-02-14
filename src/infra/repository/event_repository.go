package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"gorm.io/gorm"
)

type EventRepository struct {
	database *gorm.DB
	logger   logging.Logger
	preloads []database.PreloadEntity
}

func NewEventRepository(cfg *config.Config, preloads []database.PreloadEntity) *EventRepository {
	return &EventRepository{
		database: database.GetDb(),
		logger:   logging.NewLogger(cfg),
		preloads: preloads,
	}
}

func (r *EventRepository) Create(ctx context.Context, e model.Event) (model.Event, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.
		Create(&e).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return model.Event{}, err
	}
	tx.Commit()
	return e, nil
}

func (r *EventRepository) Update(ctx context.Context, id int, e map[string]interface{}) (model.Event, error) {
	snakeMap := map[string]interface{}{}
	for k, v := range e {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	snakeMap["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true}
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	event := new(model.Event)
	tx := r.database.WithContext(ctx).Begin()
	if err := tx.Model(event).
		Where(softDeleteWithIdExp, id).
		Updates(snakeMap).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		tx.Rollback()
		return *event, err
	}
	tx.Commit()
	return *event, nil
}

func (r *EventRepository) Delete(ctx context.Context, id int) error {
	event, err := r.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		}
		return err
	}
	if event.CreatedBy != int(ctx.Value(constant.UserIdKey).(float64)) {
		r.logger.Error(logging.Validation, logging.Permission, service_errors.PermissionDenied, nil)
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	tx := r.database.WithContext(ctx).Begin()
	model := new(model.Event)
	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	if ctx.Value(constant.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	if cnt := tx.
		Model(model).
		Where(softDeleteWithIdExp, id).
		Updates(deleteMap).
		RowsAffected; cnt == 0 {
		r.logger.Error(logging.Postgres, logging.Update, service_errors.RecordNotFound, nil)
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	tx.Commit()
	return nil
}

func (r *EventRepository) GetById(ctx context.Context, id int) (model.Event, error) {
	event := new(model.Event)

	db := r.database.
		WithContext(ctx).
		Preload("Teacher")

	err := db.
		Where(softDeleteWithIdExp, id).
		First(event).
		Error

	if err != nil {
		return *event, err
	}

	return *event, nil
}

func (r *EventRepository) GetByFilter(ctx context.Context, req filter.PaginationInput) (int64, []model.Event, error) {
	event := new(model.Event)
	var items []model.Event

	var totalRows int64 = 0
	if err := r.database.WithContext(ctx).
		Model(event).
		Where("active = ? and deleted_by is null", true).
		Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	offset := req.GetOffset()
	limit := req.GetPageSize()

	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)
	if err := db.
		Offset(offset).
		Limit(limit).
		Where("active = ? and deleted_by is null", true).
		Find(&items).Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return 0, nil, err
	}
	return totalRows, items, nil
}

func (r *EventRepository) ChangeEventStatus(ctx context.Context, id int) error {
	event := new(model.Event)
	if err := r.database.WithContext(ctx).First(event, id).Error; err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	event.Active = !event.Active

	if err := r.database.WithContext(ctx).Save(event).Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return err
	}
	return nil
}

func (r *EventRepository) ChangeCapacity(ctx context.Context, id int, capacity int) error {
	if err := r.database.WithContext(ctx).
		Model(&model.Event{}).
		Where("id = ?", id).
		Update("capacity", capacity).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return err
	}
	return nil
}
