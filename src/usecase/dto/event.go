package dto

import (
	"time"

	"github.com/farzadamr/event-manager-api/common"
	"github.com/farzadamr/event-manager-api/domain/model"
)

type CreateEvent struct {
	Title       string
	Description string
	TeacherId   int
	Capacity    int
	PosterFile  FileRef
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
	PosterFile  *FileRef
	Capacity    *int
	StartDate   *time.Time
	EndDate     *time.Time
	Location    *string
}

type EventModel struct {
	Id          int
	Title       string
	Description string
	PosterPath  string
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
	file, _ := common.TypeConverter[model.FileRef](form.PosterFile)
	return model.Event{
		Title:       form.Title,
		Description: form.Description,
		TeacherId:   form.TeacherId,
		Capacity:    form.Capacity,
		Poster:      file,
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
	if req.PosterFile != nil {
		if req.PosterFile.Path != "" {
			updates["poster_path"] = req.PosterFile.Path
		}
		if req.PosterFile.Mime != "" {
			updates["poster_mime"] = req.PosterFile.Mime
		}
		if req.PosterFile.Size != 0 {
			updates["poster_size"] = req.PosterFile.Size
		}
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
	file, _ := common.TypeConverter[FileRef](e.Poster)
	return EventModel{
		Id:          e.Id,
		Title:       e.Title,
		Description: e.Description,
		PosterPath:  file.Path,
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
