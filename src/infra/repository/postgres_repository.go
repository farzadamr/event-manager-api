package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/infra/database"
	"gorm.io/gorm"
)

const softDeleteWithIdExp = "id = ? and deleted_by is null"
const softDeleteExp = "and deleted_by is null"

type BaseRepository[TEntity any] struct {
	database *gorm.DB
}

func NewBaseRepository[TEntity any]() *BaseRepository[TEntity] {
	return &BaseRepository[TEntity]{
		database: database.GetDb(),
	}
}

func (r *BaseRepository[TEntity]) Create(ctx context.Context, entity TEntity) (TEntity, error) {
	tx := r.database.WithContext(ctx).Begin()

	err := tx.
		Create(entity).Error
	if err != nil {
		tx.Rollback()
		// log
		return entity, err
	}
	tx.Commit()
	return entity, nil
}
func (r *BaseRepository[TEntity]) Update(ctx context.Context, id int, entity map[string]interface{}) (TEntity, error) {
	snakeMap := map[string]interface{}{}
	for k, v := range entity {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	snakeMap["modified_by"] = &sql.NullInt64{Valid: true, Int64: int64(ctx.Value(constant.UserIdKey).(float64))}
	model := new(TEntity)
	tx := r.database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where(softDeleteWithIdExp, id).
		Updates(snakeMap).Error; err != nil {
		tx.Rollback()
		// log
		return *model, err
	}
	tx.Commit()
	return *model, nil
}

func (r *BaseRepository[TEntity]) Delete(ctx context.Context, id int) error {
	tx := r.database.WithContext(ctx).Begin()
	model := new(TEntity)

	deleteMap := map[string]interface{}{
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
		"deleted_by": &sql.NullInt64{Valid: true, Int64: int64(ctx.Value(constant.UserIdKey).(float64))},
	}
	if ctx.Value(constant.UserIdKey) == nil {
		// TODO : Service Error
	}
	if cnt := tx.
		Model(model).
		Where(softDeleteWithIdExp, id).
		Updates(deleteMap).RowsAffected; cnt == 0 {
		tx.Rollback()
		// log
		//TODO : Service Error
		return gorm.ErrRecordNotFound
	}
	tx.Commit()
	return nil
}

func (r *BaseRepository[TEntity]) GetById(ctx context.Context, id int) (TEntity, error) {
	model := new(TEntity)
	err := r.database.WithContext(ctx).
		Where(softDeleteWithIdExp, id).
		First(model).Error
	if err != nil {
		// log
		return *model, err
	}
	return *model, nil
}
