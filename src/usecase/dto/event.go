package dto

import (
	"time"

	"github.com/farzadamr/event-manager-api/domain/model"
)

type CreateEvent struct {
	Title       string
	Description string
	TeacherId   int
	Capacity    int
	StartDate   time.Time
	EndDate     time.Time
	Location    string
	Price       float64
}

type UpdateEvent struct {
	Id          int
	Title       *string
	Description *string
	TeacherId   *int
	Capacity    *int
	StartDate   *time.Time
	EndDate     *time.Time
	Location    *string
}

type EventModel struct {
	Id          int
	Title       string
	Description string
	Poster_Path string
	Teacher     TeacherModel
	Capacity    int
	StartDate   time.Time
	EndDate     time.Time
	Location    string
	Price       float64
	Active      bool
}

type TeacherModel struct {
	Id        int
	FirstName string
	LastName  string
}

func CreateEventToEventModel(form CreateEvent) model.Event {
	return model.Event{
		Title:       form.Title,
		Description: form.Description,
		TeacherId:   form.TeacherId,
		Capacity:    form.Capacity,
		StartDate:   form.StartDate,
		EndDate:     form.EndDate,
		Location:    form.Location,
		Price:       form.Price,
	}
}

func (req *UpdateEvent) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["Title"] = *req.Title
	}
	if req.Description != nil {
		updates["Description"] = *req.Description
	}
	if req.TeacherId != nil {
		updates["TeacherId"] = *req.TeacherId
	}
	if req.Capacity != nil {
		updates["Capacity"] = *req.Capacity
	}
	if req.StartDate != nil {
		updates["StartDate"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["EndDate"] = *req.EndDate
	}
	if req.Location != nil {
		updates["Location"] = *req.Location
	}

	return updates
}

func ToEventModel(e model.Event) EventModel {
	return EventModel{
		Id:          e.Id,
		Title:       e.Title,
		Description: e.Description,
		Poster_Path: e.Poster.Path,
		Teacher: TeacherModel{
			Id:        e.Teacher.Id,
			FirstName: e.Teacher.FirstName,
			LastName:  e.Teacher.LastName,
		},
		Capacity:  e.Capacity,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
		Location:  e.Location,
		Price:     e.Price,
		Active:    e.Active,
	}
}

func ToEventModelList(events []model.Event) []EventModel {
	if events == nil {
		return nil
	}
	result := make([]EventModel, len(events))
	for i, e := range events {
		result[i] = ToEventModel(e)
	}
	return result
}
