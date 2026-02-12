package repository

import (
	"context"

	"github.com/farzadamr/event-manager-api/domain/filter"
	"github.com/farzadamr/event-manager-api/domain/model"
)

type UserRepository interface {
	ExistsStudentNumber(ctx context.Context, studentNumber string) (bool, error)
	ExistsEmail(ctx context.Context, email string) (bool, error)
	FetchUserInfo(ctx context.Context, studentNumber string, password string) (model.User, error)
	FetchUserInfoById(ctx context.Context, id int) (model.User, error)
	GetDefaultRole(ctx context.Context) (roleId int, err error)
	CreateUser(ctx context.Context, u model.User) (model.User, error)
	GetByRoleNameByFilter(ctx context.Context, roleName string, req filter.PaginationInput) (int64, []model.User, error)
	Update(ctx context.Context, id int, e *map[string]interface{}) (model.User, error)
	AddRolesToUser(ctx context.Context, userId int, roles []int) (model.User, error)
	RemoveRolesFromUser(ctx context.Context, userId int, roles []int) error
}

type RoleRepository interface {
	Create(ctx context.Context, e model.Role) error
	Update(ctx context.Context, id int, displayName string) error
	GetById(ctx context.Context, id int) (*model.Role, error)
	GetAll(ctx context.Context) ([]model.Role, error)
}

type EventRepository interface {
	Create(ctx context.Context, e model.Event) (model.Event, error)
	Update(ctx context.Context, id int, e map[string]interface{}) (model.Event, error)
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (model.Event, error)
	GetByFilter(ctx context.Context, req filter.PaginationInput) (int64, []model.Event, error)
	ChangeEventStatus(ctx context.Context, id int) error
	ChangeCapacity(ctx context.Context, id int, capacity int) error
}

type RegistrationRepository interface {
	Create(ctx context.Context, re model.Registration) error
	FindByEventIDAndUserID(ctx context.Context, eventID, userID int) (model.Registration, error)
	FindById(ctx context.Context, id int) (model.Registration, error)
	ListByEventID(ctx context.Context, eventID int, pagination filter.PaginationInput) (int64, []model.Registration, error)
	ListByUserID(ctx context.Context, userId int, pagination filter.PaginationInput) (int64, []model.Registration, error)
	CancelByUser(ctx context.Context, eventID, userID int) error
	CancelByEvent(ctx context.Context, eventID int) error
	UpdateAttendanceList(ctx context.Context, attendanceList []model.AttendanceList) error
	GetAllAttendedByEventId(ctx context.Context, eventID int) ([]model.Registration, error)
}

type CertificateRepository interface {
	Create(ctx context.Context, r model.Certificate) (*model.Certificate, error)
	BulkCreate(ctx context.Context, certs []model.Certificate) ([]model.Certificate, error)
	MarkAsIssued(ctx context.Context, id int, file *model.FileRef, metadata *model.CertificateMetadata) error
	GetById(ctx context.Context, id int) (model.Certificate, error)
	GetByFilter(ctx context.Context, eventId int, req filter.PaginationInput) (int64, []model.Certificate, error)
	GetAllByEventId(ctx context.Context, eventId int) ([]model.Certificate, error)
	GetByRegistrationId(ctx context.Context, registrationId int) (*model.Certificate, error)
	GetByUserIdByFilter(ctx context.Context, userId int, req filter.PaginationInput) (int64, []model.Certificate, error)
	VerifyCertificate(ctx context.Context, trackingCode string) (model.Certificate, error)
}
