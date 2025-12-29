package migration

import (
	"log"

	"github.com/farzadamr/event-manager-api/constant"
	"github.com/farzadamr/event-manager-api/domain/model"
	"github.com/farzadamr/event-manager-api/infra/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Up_1() {
	database := database.GetDb()

	createTables(database)
	createDefaultUserInformation(database)
}

func Down_1() {
	// No down migration
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}
	//User
	tables = addNewTable(database, model.User{}, tables)
	tables = addNewTable(database, model.Role{}, tables)
	tables = addNewTable(database, model.UserRole{}, tables)
	//Event Management
	tables = addNewTable(database, model.Event{}, tables)
	tables = addNewTable(database, model.Registration{}, tables)
	tables = addNewTable(database, model.Certificate{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		log.Printf("Error creating tables: %v", err)
	}
	log.Println("Tables created successfully")
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createDefaultUserInformation(database *gorm.DB) {
	adminRole := model.Role{Name: constant.AdminRoleName, Display_Name: constant.AdminRoleDisplayName}
	createRoleIfNotExists(database, &adminRole)

	defaultRole := model.Role{Name: constant.DefaultRoleName, Display_Name: constant.DefaultRoleDisplayName}
	createRoleIfNotExists(database, &defaultRole)

	u := model.User{Username: constant.DefaultUserName, FirstName: "Test", LastName: "Test", Student_Number: "4010000000",
		Phone: "09120000000", Email: "test@example.com"}
	pass := "12345678"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	createAdminUserIfNotExists(database, &u, adminRole.Id)
}

func createAdminUserIfNotExists(database *gorm.DB, u *model.User, roleId int) {
	exists := 0
	database.
		Model(&model.User{}).
		Select("1").
		Where("username = ?", u.Username).
		First(&exists)
	if exists == 0 {
		database.Create(u)
		ur := model.UserRole{UserId: u.Id, RoleId: roleId}
		database.Create(&ur)
	}
}
func createRoleIfNotExists(database *gorm.DB, r *model.Role) {
	exists := 0
	database.
		Model(&model.Role{}).
		Select("1").
		Where("name = ?", r.Name).
		First(&exists)

	if exists == 0 {
		database.Create(r)
	}
}
