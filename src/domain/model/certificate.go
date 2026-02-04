package model

import (
	"database/sql"
)

type CertificateStatus string

const (
	StatusIssued  CertificateStatus = "ISSUED"
	StatusPending CertificateStatus = "CREATED"
)

type Certificate struct {
	BaseModel
	Registration   Registration `gorm:"foreignKey:RegistrationId;Constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	RegistrationId int
	IssuedAt       sql.NullTime      `gorm:"type:TIMESTAMP with time zone; null"`
	Pdf            FileRef           `gorm:"embedded;embeddedPrefix:file_"`
	Status         CertificateStatus `gorm:"type:string;default:'CREATED'"`
}
