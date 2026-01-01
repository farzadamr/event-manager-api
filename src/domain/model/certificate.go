package model

import "time"

type CertificateStatus string

const (
	StatusIssued  CertificateStatus = "ISSUED"
	StatusPending CertificateStatus = "PENDING"
)

type Certificate struct {
	BaseModel
	Registration   Registration `gorm:"foreignKey:RegistrationId;Constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	RegistrationId int
	IssuedAt       time.Time         `gorm:"type:TIMESTAMP with time zone;"`
	Pdf            FileRef           `gorm:"embedded;embeddedPrefix:file_"`
	Sent_Email     bool              `gorm:"type:boolean;default:false"`
	Status         CertificateStatus `gorm:"type:string;default:'PENDING'"`
}
