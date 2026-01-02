package dependency

import (
	contractRepository "github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/infra/database"
	infraRepository "github.com/farzadamr/event-manager-api/infra/repository"
)

func GetUserRepository() contractRepository.UserRepository {
	return infraRepository.NewUserRepository()
}

func GetEventRepository() contractRepository.EventRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{{Entity: "Teacher"}}
	return infraRepository.NewEventRepository(preloads)
}
