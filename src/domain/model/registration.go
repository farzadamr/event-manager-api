package model

import "time"

type RegistrationStatus string
type AttendanceStatus string

const (
	StatusRegistered       RegistrationStatus = "REGISTERED"
	StatusCancelledByUser  RegistrationStatus = "CANCELLED_BY_USER"
	StatusCancelledByEvent RegistrationStatus = "CANCELLED_BY_EVENT"
)
const (
	Present      AttendanceStatus = "PRESENT"
	Absent       AttendanceStatus = "ABSENT"
	NotCheckedIn AttendanceStatus = "NOT_CHECKED_IN"
)

type Registration struct {
	BaseModel
	User             User  `gorm:"foreignKey:UserId;constraint: OnUpdate:NO ACTION;OnDelete:NO ACTION;not null"`
	Event            Event `gorm:"foreignKey:EventId;constraint: OnUpdate:NO ACTION;OnDelete:NO ACTION;not null"`
	UserId           int
	EventId          int
	RegistratedAt    time.Time          `gorm:"type:TIMESTAMP with time zone;not null"`
	Status           RegistrationStatus `gorm:"type:text;default:'REGISTERED';not null"`
	AttendanceStatus AttendanceStatus   `gorm:"type:text;default:'NOT_CHECKED_IN';not null"`
}
