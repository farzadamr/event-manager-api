package repository

import (
	"context"

	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
)

type UserRepository interface {
	ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error)
	ExistsStudentNumber(ctx context.Context, studentNumber string) (bool, error)
	ExistsEmail(ctx context.Context, email string) (bool, error)
	FetchUserInfo(ctx context.Context, studentNumber string, password string) (model.User, error)
	FetchUserInfoById(ctx context.Context, id int) (model.User, error)
	GetDefaultRole(ctx context.Context) (roleId int, err error)
	CreateUser(ctx context.Context, u model.User) (model.User, error)
}

type EventRepository interface {
	Create(ctx context.Context, e model.Event) (model.Event, error)
	Update(ctx context.Context, id int, e map[string]interface{}) (model.Event, error)
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (model.Event, error)
	GetByFilter(ctx context.Context, req filter.PaginationInput) (int64, []model.Event, error)
	ChangeEventStatus(ctx context.Context, id int) error
}

type RegisterationRepository interface {
	Create(ctx context.Context, r model.Registration) error
	FindByEventIDAndUserID(ctx context.Context, eventID, userID int) (model.Registration, error)
	ListByEventID(ctx context.Context, eventID int, pagination filter.PaginationInput) (int64, []model.Registration, error)
	ListByUserID(ctx context.Context, userId int, pagination filter.PaginationInput) (int64, []model.Registration, error)
	CancelByUser(ctx context.Context, eventID, userID int) error
	CancelByEvent(ctx context.Context, eventID int) error
	UpdateAttendanceStatus(ctx context.Context, registrationId int, status model.AttendanceStatus) error
}
