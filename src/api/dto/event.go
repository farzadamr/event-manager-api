package dto

import (
	"time"

	usecase "github.com/farzadamr/event-manager-api/usecase/dto"
)

// Create Request
type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required,min=5,max=64"`
	Description string    `json:"description" binding:"required,min=16,max=256"`
	TeacherId   int       `json:"teacher_id" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required"`
	Date        time.Time `json:"date" binding:"required,date"`
	Location    string    `json:"location" binding:"required"`
	Price       float64   `json:"price"`
}

func (f CreateEventRequest) ToCreateEvent() usecase.CreateEvent {
	return usecase.CreateEvent{
		Title:       f.Title,
		Description: f.Description,
		TeacherId:   f.TeacherId,
		Capacity:    f.Capacity,
		Date:        f.Date,
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
	Date        *time.Time `json:"date,omitempty" binding:"omitempty,date"`
	Location    *string    `json:"location,omitempty" binding:"omitempty,min=5,max=64"`
}

func (f UpdateEventRequest) ToUpdateEvent() usecase.UpdateEvent {
	return usecase.UpdateEvent{
		Id:          f.Id,
		Title:       f.Title,
		Description: f.Description,
		TeacherId:   f.TeacherId,
		Capacity:    f.Capacity,
		Date:        f.Date,
		Location:    f.Location,
	}
}
