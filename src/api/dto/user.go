package dto

import usecase "github.com/farzadamr/event-manager-api/usecase/dto"

type RegisterUserByStudentNumberRequest struct {
	FirstName     string `json:"firstName" binding:"required,min=3"`
	LastName      string `json:"lastName" binding:"required,min=3"`
	StudentNumber string `json:"studentNumber" binding:"required,min=10,max=10"`
	Email         string `json:"email" binding:"min=6"`
	Password      string `json:"password" binding:"required,min=6"`
}

func (from RegisterUserByStudentNumberRequest) ToRegisterUserByStudentNumber() usecase.RegisterByStudentNumber {
	return usecase.RegisterByStudentNumber{
		FirstName:     from.FirstName,
		LastName:      from.LastName,
		StudentNumber: from.StudentNumber,
		Email:         from.Email,
		Password:      from.Password,
	}
}
