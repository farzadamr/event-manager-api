package usecase

import (
	"context"
	"errors"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"github.com/farzadamr/event-manager-api/usecase/dto"
	"gorm.io/gorm"
)

type EventUsecase struct {
	cfg             *config.Config
	eventRepository repository.EventRepository
	userRepository  repository.UserRepository
}

func NewEventUsecase(cfg *config.Config, eventRepo repository.EventRepository, userRepo repository.UserRepository) *EventUsecase {
	return &EventUsecase{cfg: cfg, eventRepository: eventRepo, userRepository: userRepo}
}

func (u *EventUsecase) PublishEvent(ctx context.Context, req dto.CreateEvent) error {
	_, err := u.userRepository.FetchUserInfoById(ctx, req.TeacherId)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}

	_, err = u.eventRepository.Create(ctx, dto.CreateEventToEventModel(req))
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}

	return nil
}

func (u *EventUsecase) Update(ctx context.Context, req dto.UpdateEvent) (model.Event, error) {
	event, err := u.eventRepository.GetById(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Event{}, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		}
		return model.Event{}, err
	}

	userId, _ := ctx.Value(constant.UserIdKey).(float64)
	if event.CreatedBy != int(userId) {
		return model.Event{}, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	updates := req.ToUpdateMap()
	if len(updates) == 0 {
		return model.Event{}, errors.New("no fields to update")
	}

	res, err := u.eventRepository.Update(ctx, req.Id, updates)
	if err != nil {
		return model.Event{}, err
	}
	return res, nil
}

func (u *EventUsecase) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}

	if err := u.eventRepository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *EventUsecase) ChangeEventStatus(ctx context.Context, id int) error {
	if id == 0 {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
	}
	if err := u.eventRepository.ChangeEventStatus(ctx, id); err != nil {
		return err
	}
	return nil
}

func (u *EventUsecase) GetByFilter(ctx context.Context, req filter.PaginationInput) (*filter.PagedList[dto.EventModel], error) {
	count, events, err := u.eventRepository.GetByFilter(ctx, req)
	if err != nil {
		return nil, err
	}
	dtoEvents := dto.ToEventModelList(events)
	return filter.NewPagedList(&dtoEvents, count, req.PageNumber, int64(req.PageSize)), nil
}

func (u *EventUsecase) GetById(ctx context.Context, id int) (dto.EventModel, error) {
	event, err := u.eventRepository.GetById(ctx, id)
	if err != nil {
		return dto.EventModel{}, err
	}

	result := dto.ToEventModel(event)
	return result, nil
}
