package model

type User struct {
	BaseModel
	StudentNumber string  `gorm:"type:string;size:10;unique;not null"`
	FirstName     string  `gorm:"type:string;size:15;not null"`
	LastName      string  `gorm:"type:string;size:25;not null"`
	EnglishName   *string `gorm:"type:string;size:15;null"`
	Phone         *string `gorm:"type:string;size:11;null"`
	Email         string  `gorm:"type:string;size:64;unique"`
	Password      string  `gorm:"type:string;size:64;not null"`
	Active        bool    `gorm:"default:true"`
	UserRoles     []UserRole
}

type Role struct {
	BaseModel
	Name        string `gorm:"type:string;size:25;not null;unique"`
	DisplayName string `gorm:"type:string;size:25;not null;unique"`
	UserRoles   []UserRole
}

type UserRole struct {
	BaseModel
	UserId int
	RoleId int
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
