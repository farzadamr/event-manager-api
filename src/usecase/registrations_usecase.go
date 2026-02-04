package usecase

import (
	"context"

	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/pkg/service_errors"
	"github.com/farzadamr/event-manager-api/usecase/dto"
)

type RegistrationUseCase struct {
	config       *config.Config
	registerRepo repository.RegistrationRepository
	userRepo     repository.UserRepository
}

func NewRegistrationUseCase(cfg *config.Config, RegisterRepo repository.RegistrationRepository, userRepo repository.UserRepository) *RegistrationUseCase {
	return &RegistrationUseCase{config: cfg, registerRepo: RegisterRepo, userRepo: userRepo}
}

// TODO: implement ownership checker for teacher role => teachers can only access their own events to update and get attendance list.

func (u *RegistrationUseCase) GetRegistrations(ctx context.Context, eventId int, req filter.PaginationInput) (*filter.PagedList[dto.Registration], error) {
	count, registrations, err := u.registerRepo.ListByEventID(ctx, eventId, req)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	dtoRegistrations := dto.ToRegistrationList(registrations)
	return filter.NewPagedList(&dtoRegistrations, count, req.PageNumber, int64(req.PageSize)), nil
}

func (u *RegistrationUseCase) AttendanceList(ctx context.Context, eventId int, req filter.PaginationInput) (*filter.PagedList[dto.Attendance], error) {
	registrations, err := u.GetRegistrations(ctx, eventId, req)
	if err != nil {
		return nil, err
	}
	attendanceList := dto.ToAttendanceList(registrations.Items)
	count := int64(len(attendanceList))
	return filter.NewPagedList(&attendanceList, count, req.PageNumber, int64(req.PageSize)), nil
}

func (u *RegistrationUseCase) UpdateAttendanceList(ctx context.Context, attendanceList []dto.AttendanceRequest) error {
	userId := int(ctx.Value(constant.UserIdKey).(float64))
	if userId == 0 {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	user, err := u.userRepo.FetchUserInfoById(ctx, userId)

	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	hasValidRole := false
	for _, userRole := range user.UserRoles {
		if userRole.Role.Name == constant.AdminRoleName {
			hasValidRole = true
		}
	}
	if !hasValidRole {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	if err = u.registerRepo.UpdateAttendanceList(ctx, dto.ToAttendanceListModel(attendanceList)); err != nil {
		return err
	}
	return nil
}
