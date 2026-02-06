package model

import (
	"time"
)

type Certificate struct {
	BaseModel
	Registration   Registration `gorm:"foreignKey:RegistrationId;Constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	RegistrationId int          `gorm:"uniqueIndex"`
	IssuedAt       *time.Time   `gorm:"type:TIMESTAMP with time zone; null"`
	Pdf            FileRef      `gorm:"embedded;embeddedPrefix:file_"`
	TrackingCode   string       `gorm:"type:varchar(64);uniqueIndex;not null"`
}
