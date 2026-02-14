package repository

import (
	"context"
	"errors"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"gorm.io/gorm"
)

type RoleRepository struct {
	database *gorm.DB
	logger   logging.Logger
}

func NewRoleRepository(cfg *config.Config) *RoleRepository {
	return &RoleRepository{
		database: database.GetDb(),
		logger:   logging.NewLogger(cfg),
	}
}

func (r *RoleRepository) Create(ctx context.Context, e model.Role) error {
	err := r.database.WithContext(ctx).Create(&e).Error
	if err != nil {
		r.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return err
	}
	return nil
}

func (r *RoleRepository) Update(ctx context.Context, id int, displayName string) error {
	db := r.database.WithContext(ctx)
	res := db.Model(&model.Role{}).
		Where("id = ?", id).
		Updates(model.Role{Name: displayName})
	if res.Error != nil {
		r.logger.Error(logging.Postgres, logging.Update, res.Error.Error(), nil)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *RoleRepository) GetById(ctx context.Context, id int) (*model.Role, error) {
	role := model.Role{}
	db := r.database.WithContext(ctx)
	err := db.First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		}
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}
	return &role, nil
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	db := r.database.WithContext(ctx)
	err := db.Find(&roles).Error
	if err != nil {
		return roles, err
	}
	return roles, nil
}
