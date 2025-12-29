package dependency

import (
	contractRepository "github.com/farzadamr/event-manager-api/domain/repository"
	infraRepository "github.com/farzadamr/event-manager-api/infra/repository"
)

func GetUserRepository() contractRepository.UserRepository {
	return infraRepository.NewUserRepository()
}
