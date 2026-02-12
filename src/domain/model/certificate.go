package model

import (
	"time"
)

type IssueStatus string

const Pending IssueStatus = "PENDING"
const Issued IssueStatus = "ISSUED"

type Certificate struct {
	BaseModel
	Registration   Registration `gorm:"foreignKey:RegistrationId;Constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	RegistrationId int          `gorm:"uniqueIndex"`
	EventId        int
	IssuedAt       *time.Time           `gorm:"type:TIMESTAMP with time zone; null"`
	Pdf            *FileRef             `gorm:"embedded;embeddedPrefix:file_"`
	TrackingCode   string               `gorm:"type:varchar(64);uniqueIndex;not null"`
	Status         IssueStatus          `gorm:"default:'PENDING';not null"`
	Metadata       *CertificateMetadata `gorm:"embedded;embeddedPrefix:metadata_"`
}
type CertificateMetadata struct {
	UserName    string  `gorm:"type:varchar(32)"`
	EnglishName *string `gorm:"type:varchar(64)"`
	Date        string  `gorm:"type:varchar(10)"`
	Duration    string  `gorm:"type:varchar(4)"`
	EventName   string  `gorm:"type:varchar(128)"`
}
