package dto

import (
	"fmt"
	"time"

	"github.com/farzadamr/event-manager-api/domain/model"
)

type Certificate struct {
	Id           int
	RegisterId   int
	Name         string
	PdfPath      string
	IssuedAt     *time.Time
	TrackingCode string
	Status       string
}

func ToCertificateList(certificates []model.Certificate) []Certificate {
	if len(certificates) == 0 {
		return nil
	}
	items := make([]Certificate, len(certificates))
	for i, c := range certificates {
		items[i] = ToCertificateDto(c)
	}
	return items
}

func ToCertificateDto(model model.Certificate) Certificate {
	return Certificate{
		Id:           model.Id,
		RegisterId:   model.RegistrationId,
		Name:         model.Registration.User.FirstName + " " + model.Registration.User.LastName,
		PdfPath:      model.Pdf.Path,
		IssuedAt:     model.IssuedAt,
		TrackingCode: model.TrackingCode,
		Status:       string(model.Status),
	}
}

type UserCertificate struct {
	Id           int
	EventName    string
	TrackingCode string
	Status       string
	downloadLink string
}

func ToUserCertificateList(certificates []model.Certificate) []UserCertificate {
	if len(certificates) == 0 {
		return nil
	}
	items := make([]UserCertificate, len(certificates))
	for i, c := range certificates {
		items[i] = ToUserCertificate(c)
	}
	return items
}

func ToUserCertificate(model model.Certificate) UserCertificate {
	downloadLink := fmt.Sprintf("http://localhost:5005/api/v1/certificates/%d/download", model.Id)
	return UserCertificate{
		Id:           model.Id,
		EventName:    model.Registration.Event.Title,
		TrackingCode: model.TrackingCode,
		Status:       string(model.Status),
		downloadLink: downloadLink,
	}
}
