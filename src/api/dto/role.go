package dto

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=16"`
	DisplayName string `json:"display_name" binding:"required,min=3,max=16"`
}

type UpdateRoleRequest struct {
	Id          int    `json:"id" binding:"required"`
	DisplayName string `json:"display_name" binding:"required,min=3,max=16"`
}
