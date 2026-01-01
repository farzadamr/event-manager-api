package dto

import "time"

type CreateEvent struct {
	Title       string
	Description string
	TeacherId   int
	Capacity    int
	Date        time.Time
	Location    string
	Price       float64
}

type UpdateEvent struct {
	Title       string
	Description string
	TeacherId   int
	Capacity    int
	Date        time.Time
	Location    string
}

type EventModel struct {
	Id          int
	Title       string
	Description string
	Poster_Path string
	Teacher     TeacherModel
	Capacity    int
	Date        time.Time
	Location    string
	Price       float64
	Active      bool
}

type TeacherModel struct {
	Id        int
	FirstName string
	LastName  string
}
