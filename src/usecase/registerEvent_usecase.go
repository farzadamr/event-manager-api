package usecase

import (
	"context"
	"time"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
)

type RegisterEventUsecase struct {
	cfg                *config.Config
	logger             logging.Logger
	registerRepository repository.RegistrationRepository
	eventRepository    repository.EventRepository
	userRepository     repository.UserRepository
}

func NewRegisterEventUsecase(cfg *config.Config,
	registerRepo repository.RegistrationRepository,
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository) *RegisterEventUsecase {
	return &RegisterEventUsecase{
		cfg:                cfg,
		logger:             logging.NewLogger(cfg),
		registerRepository: registerRepo,
		userRepository:     userRepo,
		eventRepository:    eventRepo,
	}
}

func (u *RegisterEventUsecase) RegisterForEvent(ctx context.Context, eventID, userId int) error {
	register, err := u.registerRepository.FindByEventIDAndUserID(ctx, eventID, userId)
	if err == nil {
		return &service_errors.ServiceError{EndUserMessage: "registration already exist"}
	}
	if register.Event.Capacity == 0 {
		u.logger.Error(logging.Internal, logging.NoCapacity, "", nil)
		return &service_errors.ServiceError{EndUserMessage: service_errors.NoCapacity}
	}
	registration := model.Registration{
		EventId:       eventID,
		UserId:        userId,
		RegistratedAt: time.Now().UTC(),
	}
	err = u.registerRepository.Create(ctx, registration)
	if err != nil {
		return err
	}
	//decrease capacity
	newCapacity := register.Event.Capacity - 1
	if err = u.eventRepository.ChangeCapacity(ctx, eventID, newCapacity); err != nil {
		return err
	}
	return nil
}

func (u *RegisterEventUsecase) CancelRegistration(ctx context.Context, eventId, userId int) error {
	if userId != int(ctx.Value(constant.UserIdKey).(float64)) {
		u.logger.Error(logging.Validation, logging.Permission, "", nil)
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	err := u.registerRepository.CancelByUser(ctx, eventId, userId)
	if err != nil {
		return err
	}
	return nil
}
