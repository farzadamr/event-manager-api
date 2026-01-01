package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"gorm.io/gorm"
)

type EventRepository struct {
	database *gorm.DB
	preloads []database.PreloadEntity
}

func NewEventRepository(preloads []database.PreloadEntity) *EventRepository {
	return &EventRepository{database: database.GetDb(), preloads: preloads}
}

func (r *EventRepository) Create(ctx context.Context, e model.Event) (model.Event, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.
		Create(e).Error
	if err != nil {
		tx.Rollback()
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
	model := new(model.Event)
	tx := r.database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where(softDeleteWithIdExp, id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		return *model, err
	}
	tx.Commit()
	return *model, nil
}

func (r *EventRepository) Delete(ctx context.Context, id int) error {
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

func (r *EventRepository) GetByFilter(ctx context.Context, req filter.PaginationInput) (int64, *[]model.Event, error) {
	event := new(model.Event)
	var items *[]model.Event
	db := database.Preload(r.database, r.preloads)

	var totalRows int64 = 0
	if err := db.
		Model(event).
		Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	offset := req.GetOffset()
	limit := req.GetPageSize()

	if err := db.Offset(offset).Limit(limit).Find(items).Error; err != nil {
		return 0, nil, err
	}
	return totalRows, items, nil
}

func (r *EventRepository) ChangeEventStatus(ctx context.Context, id int) error {
	event := new(model.Event)
	if err := r.database.WithContext(ctx).First(event, id).Error; err != nil {
		return err //not found
	}

	event.Active = !event.Active

	if err := r.database.WithContext(ctx).Save(event).Error; err != nil {
		return err
	}
	return nil
}
