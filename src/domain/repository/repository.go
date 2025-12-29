package repository

import (
	"context"

	"github.com/farzadamr/event-manager-api/domain/model"
)

type BaseRepository[TEntity any] interface {
	Create(ctx context.Context, entity TEntity) (TEntity, error)
	Update(ctx context.Context, id int, entity map[string]interface{}) (TEntity, error)
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (TEntity, error)
}

type UserRepository interface {
	ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error)
	ExistsStudentNumber(ctx context.Context, studentNumber string) (bool, error)
	ExistsEmail(ctx context.Context, email string) (bool, error)
	FetchUserInfo(ctx context.Context, studentNumber string, password string) (model.User, error)
	GetDefaultRole(ctx context.Context) (roleId int, err error)
	CreateUser(ctx context.Context, u model.User) (model.User, error)
}

type RoleRepository interface {
	BaseRepository[model.Role]
}

type CertificateRepository interface {
	BaseRepository[model.Certificate]
}

type RegistrationRepository interface {
	BaseRepository[model.Registration]
}

type EventRepository interface {
	BaseRepository[model.Event]
}
