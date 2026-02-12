package dto

import "github.com/farzadamr/event-manager-api/domain/model"

type TokenDetail struct {
	AccessToken            string
	RefreshToken           string
	AccessTokenExpireTime  int64
	RefreshTokenExpireTime int64
}

type RegisterByStudentNumber struct {
	FirstName     string
	LastName      string
	StudentNumber string
	EnglishName   *string
	PhoneNumber   *string
	Email         string
	Password      string
}

func ToUserModel(from RegisterByStudentNumber) model.User {
	return model.User{
		FirstName:     from.FirstName,
		LastName:      from.LastName,
		StudentNumber: from.StudentNumber,
		EnglishName:   from.EnglishName,
		Phone:         from.PhoneNumber,
		Email:         from.Email,
	}
}

type UserDto struct {
	Id            int
	StudentNumber string
	Name          string
	EnglishName   *string
	PhoneNumber   *string
	Email         string
	Active        bool
	Roles         []string
}

func ToUserDto(model model.User) UserDto {
	dto := UserDto{
		Id:            model.Id,
		StudentNumber: model.StudentNumber,
		Name:          model.FirstName + " " + model.LastName,
		EnglishName:   model.EnglishName,
		PhoneNumber:   model.Phone,
		Email:         model.Email,
		Active:        model.Active,
	}
	var roles []string
	for _, role := range model.UserRoles {
		roles = append(roles, role.Role.DisplayName)
	}
	dto.Roles = roles
	return dto
}

func ToUserDtoList(model *[]model.User) *[]UserDto {
	dtoList := make([]UserDto, len(*model))
	for i, user := range *model {
		dtoList[i] = ToUserDto(user)
	}
	return &dtoList
}

type UpdateUserDto struct {
	Id          int
	FirstName   *string
	LastName    *string
	Email       *string
	EnglishName *string
	Phone       *string
}

func (req *UpdateUserDto) ToMap() *map[string]interface{} {
	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.EnglishName != nil {
		updates["english_name"] = *req.EnglishName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	return &updates
}
