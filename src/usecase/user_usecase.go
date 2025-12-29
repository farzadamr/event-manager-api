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
	cfg          *config.Config
	repository   repository.UserRepository
	tokenUsecase *TokenUsecase
}

func NewUserUsecase(cfg *config.Config, repository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		cfg:          cfg,
		repository:   repository,
		tokenUsecase: NewTokenUsecase(cfg),
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

func (u *UserUsecase) LoginByStudentnumber(ctx context.Context, studenNumber string, password string) (*dto.TokenDetail, error) {
	user, err := u.repository.FetchUserInfo(ctx, studenNumber, password)
	if err != nil {
		return nil, err
	}

	tokenDto := tokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, StudentNumber: user.Student_Number}
	if len(user.UserRoles) > 0 {
		for _, ur := range user.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := u.tokenUsecase.GenerateToken(tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}
