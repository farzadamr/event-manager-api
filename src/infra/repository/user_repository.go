package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	userFilterExp  string = "student_number = ?"
	countFilterExp string = "count(*) > 0"
)

type UserRepository struct {
	database *gorm.DB
	logger   logging.Logger
	preloads []database.PreloadEntity
}

func NewUserRepository(cfg *config.Config, preloads []database.PreloadEntity) *UserRepository {
	return &UserRepository{database: database.GetDb(),
		logger:   logging.NewLogger(cfg),
		preloads: preloads,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	roleId, err := r.GetDefaultRole(ctx)
	if err != nil {
		return u, err
	}
	tx := r.database.WithContext(ctx).Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return u, err
	}
	err = tx.Create(&model.UserRole{UserId: u.Id, RoleId: roleId}).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return u, err
	}
	tx.Commit()
	return u, nil
}

func (r *UserRepository) FetchUserInfo(ctx context.Context, username string, password string) (model.User, error) {
	var user model.User
	err := r.database.WithContext(ctx).
		Model(&model.User{}).
		Where(userFilterExp, username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error

	if err != nil {
		return user, &service_errors.ServiceError{EndUserMessage: service_errors.UsernameOrPasswordInvalid}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, &service_errors.ServiceError{EndUserMessage: service_errors.UsernameOrPasswordInvalid}
	}

	return user, nil
}

func (r *UserRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	if err := r.database.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) ExistsStudentNumber(ctx context.Context, studentNumber string) (bool, error) {
	var exists bool
	if err := r.database.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where(userFilterExp, studentNumber).
		Find(&exists).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetDefaultRole(ctx context.Context) (roleId int, err error) {

	if err = r.database.WithContext(ctx).Model(&model.Role{}).
		Select("id").
		Where("name = ?", constant.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}

func (r *UserRepository) FetchUserInfoById(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := r.database.
		Model(&model.User{}).
		Where("id = ?", id).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		First(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByRoleNameByFilter(ctx context.Context, roleName string, req filter.PaginationInput) (int64, []model.User, error) {
	var users []model.User
	var totalCount int64

	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)

	limit := req.GetPageSize()
	offset := req.GetOffset()

	query := db.
		Model(&model.User{}).
		Joins("INNER JOIN user_roles ur ON users.id = ur.user_id").
		Joins("INNER JOIN roles r ON ur.role_id = r.id").
		Where("ur.name = ?", roleName)

	if err := query.Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	err := query.
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	if err != nil {
		return 0, users, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}

	return totalCount, users, nil
}

func (r *UserRepository) Update(ctx context.Context, id int, e *map[string]interface{}) (model.User, error) {
	snakeMap := map[string]interface{}{}
	for k, v := range *e {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	snakeMap["modified_by"] = &sql.NullInt64{Valid: true, Int64: int64(ctx.Value(constant.UserIdKey).(float64))}

	user := new(model.User)
	tx := r.database.WithContext(ctx).Begin()
	if err := tx.Model(user).
		Where(softDeleteWithIdExp, id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return *user, err
	}
	tx.Commit()
	return *user, nil
}

func (r *UserRepository) AddRolesToUser(ctx context.Context, userId int, roleIds []int) (model.User, error) {
	var user model.User

	db := r.database.WithContext(ctx)
	db = database.Preload(db, r.preloads)

	if err := db.First(&user, userId).Error; err != nil {
		return user, err
	}

	var roles []model.Role
	if err := db.Find(&roles, roleIds).Error; err != nil {
		return user, err
	}

	err := db.Model(&user).Association("Roles").Append(roles)
	if err != nil {
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return user, err
	}
	return user, nil
}
func (r *UserRepository) RemoveRolesFromUser(ctx context.Context, userId int, roleIds []int) error {
	var user model.User

	user.Id = userId

	var roles []model.Role
	for _, id := range roleIds {
		roles = append(roles, model.Role{BaseModel: model.BaseModel{Id: id}})
	}

	err := r.database.WithContext(ctx).Model(&user).Association("Roles").Delete(roles)
	if err != nil {
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return err
	}
	return nil
}
