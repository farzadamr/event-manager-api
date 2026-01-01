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
	GetDefaultRole(ctx context.Context) (roleId int, err error)
	CreateUser(ctx context.Context, u model.User) (model.User, error)
}

type EventRepository interface {
	Create(ctx context.Context, e model.Event) (model.Event, error)
	Update(ctx context.Context, id int, e map[string]interface{}) (model.Event, error)
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (model.Event, error)
	GetByFilter(ctx context.Context, req filter.PaginationInput) (int64, *[]model.Event, error)
	ChangeEventStatus(ctx context.Context, id int) error
}
