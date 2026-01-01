package model

import "time"

type Event struct {
	BaseModel
	Title       string  `gorm:"type:string;size:64;not null;unique"`
	Description string  `gorm:"type:string;size:1024;not null"`
	Poster      FileRef `gorm:"embedded;embeddedPrefix:poster_"`
	Teacher     User    `gorm:"foreignKey:TeacherId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	TeacherId   int
	Capacity    int
	Date        time.Time `gorm:"type:TIMESTAMP with time zone;not null"`
	Location    string    `gorm:"type:string;size:128;not null"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Active      bool      `gorm:"default:true"`
}
