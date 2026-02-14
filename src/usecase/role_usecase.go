package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"github.com/farzadamr/event-manager-api/usecase/dto"
	"gorm.io/gorm"
)

type RoleUsecase struct {
	cfg            *config.Config
	logger         logging.Logger
	roleRepository repository.RoleRepository
}

func NewRoleUsecase(cfg *config.Config, roleRepo repository.RoleRepository) *RoleUsecase {
	return &RoleUsecase{
		cfg:            cfg,
		logger:         logging.NewLogger(cfg),
		roleRepository: roleRepo,
	}
}

func (u *RoleUsecase) Create(ctx context.Context, req dto.CreateRole) error {
	role := model.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
	}
	err := u.roleRepository.Create(ctx, role)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &service_errors.ServiceError{EndUserMessage: service_errors.DuplicatedKey}
		}
		return &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}
	return nil
}

func (u *RoleUsecase) Update(ctx context.Context, id int, displayName string) error {
	err := u.roleRepository.Update(ctx, id, displayName)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &service_errors.ServiceError{EndUserMessage: service_errors.DuplicatedKey}
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		}
		return &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}
	return nil
}

func (u *RoleUsecase) GetById(ctx context.Context, id int) (dto.Role, error) {
	var roleDto dto.Role
	role, err := u.roleRepository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return roleDto, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		}
		return roleDto, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}

	roleDto, err = common.TypeConverter[dto.Role](role)
	if err != nil {
		log.Println(err)
		return roleDto, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}
	return roleDto, nil
}

func (u *RoleUsecase) GetAll(ctx context.Context) ([]dto.Role, error) {
	var rolesDto []dto.Role

	roles, err := u.roleRepository.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return rolesDto, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}

	for _, role := range roles {
		roleDto, convErr := common.TypeConverter[dto.Role](role)
		if convErr != nil {
			log.Println(convErr)
			return rolesDto, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		rolesDto = append(rolesDto, roleDto)
	}
	return rolesDto, nil
}
