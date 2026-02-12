package dto

import usecase "github.com/farzadamr/event-manager-api/usecase/dto"

type RegisterUserByStudentNumberRequest struct {
	FirstName     string `json:"firstName" binding:"required,min=3"`
	LastName      string `json:"lastName" binding:"required,min=3"`
	StudentNumber string `json:"student_number" binding:"required,min=10,max=10"`
	EnglishName   string `json:"english_name" binding:"min=3,max=30"`
	PhoneNumber   string `json:"phone_number" binding:"mobile,min=11,max=11"`
	Email         string `json:"email" binding:"email,min=6.max=32"`
	Password      string `json:"password" binding:"required,password,min=6,max=16"`
}
type LoginByStudentNumberRequest struct {
	StudentNumber string `json:"student_number" binding:"required,min=10,max=10"`
	Password      string `json:"password" binding:"required,password,min=6,max=16"`
}

func (from RegisterUserByStudentNumberRequest) ToRegisterUserByStudentNumber() usecase.RegisterByStudentNumber {
	return usecase.RegisterByStudentNumber{
		FirstName:     from.FirstName,
		LastName:      from.LastName,
		StudentNumber: from.StudentNumber,
		EnglishName:   &from.EnglishName,
		PhoneNumber:   &from.PhoneNumber,
		Email:         from.Email,
		Password:      from.Password,
	}
}

type UpdateUserRequest struct {
	Id          int     `json:"id" binding:"required"`
	FirstName   *string `json:"firstName" binding:"min=3,max=16"`
	LastName    *string `json:"lastName" binding:"min=3,max=16"`
	Email       *string `json:"email" binding:"email,min=6,max=32"`
	EnglishName *string `json:"englishName" binding:"min=3,max=32"`
	Phone       *string `json:"phone" binding:"mobile"`
}

type RoleToUserRequest struct {
	UserId  int   `json:"user_id" binding:"required"`
	RoleIds []int `json:"role_ids" binding:"required"`
}
