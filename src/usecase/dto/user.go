package dto

import "github.com/farzadamr/event-manager-api/domain/model"

type RegisterByStudentNumber struct {
	FirstName     string
	LastName      string
	StudentNumber string
	Email         string
	Password      string
}

func ToUserModel(from RegisterByStudentNumber) model.User {
	return model.User{
		FirstName: from.FirstName,
		LastName:  from.LastName,
		Email:     from.Email,
	}
}
