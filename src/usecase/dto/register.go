package dto

import (
	"time"

	"github.com/farzadamr/event-manager-api/domain/model"
)

type Registration struct {
	Id               int
	User             User
	Event            Event
	RegistratedAt    time.Time
	Status           string
	AttendanceStatus string
}

type User struct {
	Id            int
	Name          string
	StudentNumber string
}
type Event struct {
	Id    int
	Title string
}

func ToRegistrationDto(m model.Registration) Registration {
	return Registration{
		Id: m.Id,
		User: User{
			Id:            m.UserId,
			Name:          m.User.FirstName + " " + m.User.LastName,
			StudentNumber: m.User.Student_Number,
		},
		Event: Event{
			Id:    m.EventId,
			Title: m.Event.Title,
		},
		RegistratedAt:    m.RegistratedAt,
		Status:           string(m.Status),
		AttendanceStatus: string(m.AttendanceStatus),
	}
}
func ToRegistrationList(registrations []model.Registration) []Registration {
	if registrations == nil {
		return nil
	}
	result := make([]Registration, len(registrations))
	for i, e := range registrations {
		result[i] = ToRegistrationDto(e)
	}
	return result
}

type Attendance struct {
	Id     int
	User   User
	Status string
}
type AttendanceRequest struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

func ToAttendanceDto(in Registration) Attendance {
	return Attendance{
		Id:     in.Id,
		User:   in.User,
		Status: in.AttendanceStatus,
	}
}
func ToAttendanceList(in *[]Registration) []Attendance {
	if in == nil {
		return nil
	}
	result := make([]Attendance, len(*in))
	for k, v := range *in {
		result[k] = ToAttendanceDto(v)
	}
	return result
}

func ToAttendanceListModel(in []AttendanceRequest) []model.AttendanceList {
	if in == nil {
		return nil
	}
	var alm []model.AttendanceList
	for _, v := range in {
		alm = append(alm, model.AttendanceList{
			Id:               v.Id,
			AttendanceStatus: model.CheckAttendanceStatus(v.Status),
		})
	}
	return alm
}
