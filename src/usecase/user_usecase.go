package usecase

import (
	"context"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"github.com/farzadamr/event-manager-api/usecase/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	cfg        *config.Config
	repository repository.UserRepository
}

func NewUserUsecase(cfg *config.Config, repository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		cfg:        cfg,
		repository: repository,
	}
}

func (u *UserUsecase) RegisterByStudentNumber(ctx context.Context, req dto.RegisterByStudentNumber) error {
	user := dto.ToUserModel(req)

	exists, err := u.repository.ExistsEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
	}
	exists, err = u.repository.ExistsStudentNumber(ctx, req.StudentNumber)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.StudentNumberExists}
	}
	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		//log
		return err
	}

	user.Password = string(hp)
	_, err = u.repository.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
