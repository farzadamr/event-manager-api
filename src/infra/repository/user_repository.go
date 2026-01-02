package repository

import (
	"context"

	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
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
}

func NewUserRepository() *UserRepository {
	return &UserRepository{database: database.GetDb()}
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
		return u, err
	}
	err = tx.Create(&model.UserRole{UserId: u.Id, RoleId: roleId}).Error
	if err != nil {
		tx.Rollback()
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
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error) {
	var exists bool
	if err := r.database.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
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
		Find(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
