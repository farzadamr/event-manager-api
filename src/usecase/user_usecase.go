package usecase

import (
	"context"
	"strings"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
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

func (u *UserUsecase) LoginByStudentNumber(ctx context.Context, studentNumber string, password string) (*dto.TokenDetail, error) {
	user, err := u.repository.FetchUserInfo(ctx, studentNumber, password)
	if err != nil {
		return nil, err
	}

	tokenDto := tokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, StudentNumber: user.StudentNumber}
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

func (u *UserUsecase) GetUserListByRoleName(ctx context.Context, roleName string, req filter.PaginationInput) (*filter.PagedList[dto.UserDto], error) {
	if err := u.ValidateRoleName(roleName); err != nil {
		return nil, err
	}

	count, users, err := u.repository.GetByRoleNameByFilter(ctx, roleName, req)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}

	usersDto := dto.ToUserDtoList(&users)
	return filter.NewPagedList(usersDto, count, req.PageNumber, int64(req.PageSize)), nil
}

func (u *UserUsecase) Update(ctx context.Context, req dto.UpdateUserDto) error {
	_, err := u.permissionUpdateCheck(ctx, req.Id)
	if err != nil {
		return err
	}
	updates := req.ToMap()
	if len(*updates) == 0 {
		return &service_errors.ServiceError{EndUserMessage: "no fields to update"}
	}

	_, err = u.repository.Update(ctx, req.Id, updates)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) ValidateRoleName(roleName string) error {
	if strings.TrimSpace(roleName) == "" {
		return &service_errors.ServiceError{EndUserMessage: "role name is empty"}
	}
	var validRoles = map[string]bool{
		constant.DefaultRoleName: true,
		constant.TeacherRoleName: true,
		constant.AdminRoleName:   true,
	}
	if !validRoles[roleName] {
		return &service_errors.ServiceError{EndUserMessage: "role name is invalid"}
	}
	return nil
}

func (u *UserUsecase) permissionUpdateCheck(ctx context.Context, id int) (*model.User, error) {
	user, err := u.repository.FetchUserInfoById(ctx, id)
	if err != nil {
		return nil, err
	}
	roleCheck := false
	for _, ur := range user.UserRoles {
		if ur.Role.Name == constant.AdminRoleName {
			roleCheck = true
		}
	}
	if !roleCheck {
		requestUserId, _ := ctx.Value(constant.UserIdKey).(float64)
		if user.Id != int(requestUserId) {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
		}
		return &user, nil
	}
	return &user, nil
}
