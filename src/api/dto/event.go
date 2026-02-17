package dto

import (
	"time"

	usecase "github.com/farzadamr/event-manager-api/usecase/dto"
)

// Create Request
type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required,min=5,max=64"`
	Description string    `json:"description" binding:"required,min=16,max=1024"`
	TeacherId   int       `json:"teacher_id" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required,date"`
	EndDate     time.Time `json:"end_date" binding:"required,date"`
	Location    string    `json:"location" binding:"required"`
	Price       float64   `json:"price"`
}

func (f CreateEventRequest) ToCreateEvent() usecase.CreateEvent {
	return usecase.CreateEvent{
		Title:       f.Title,
		Description: f.Description,
		TeacherId:   f.TeacherId,
		Capacity:    f.Capacity,
		StartDate:   f.StartDate,
		EndDate:     f.EndDate,
		Location:    f.Location,
		Price:       f.Price,
	}
}

// Update Request
type UpdateEventRequest struct {
	Id          int        `json:"id"`
	Title       *string    `json:"title,omitempty" binding:"omitempty,min=5,max=64"`
	Description *string    `json:"description,omitempty" binding:"omitempty,min=16,max=256"`
	TeacherId   *int       `json:"teacher_id,omitempty" binding:"omitempty,min=1"`
	Capacity    *int       `json:"capacity,omitempty" binding:"omitempty,min=1,max=250"`
	StartDate   *time.Time `json:"start_date,omitempty" binding:"omitempty,date"`
	EndDate     *time.Time `json:"end_date,omitempty" binding:"omitempty,date"`
	Location    *string    `json:"location,omitempty" binding:"omitempty,min=5,max=64"`
}

func (f UpdateEventRequest) ToUpdateEvent() usecase.UpdateEvent {
	return usecase.UpdateEvent{
		Id:          f.Id,
		Title:       f.Title,
		Description: f.Description,
		TeacherId:   f.TeacherId,
		Capacity:    f.Capacity,
		StartDate:   f.StartDate,
		EndDate:     f.EndDate,
		Location:    f.Location,
	}
}
