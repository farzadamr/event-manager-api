package dependency

import (
	"github.com/farzadamr/event-manager-api/config"
	contractRepository "github.com/farzadamr/event-manager-api/domain/repository"
	"github.com/farzadamr/event-manager-api/infra/database"
	infraRepository "github.com/farzadamr/event-manager-api/infra/repository"
)

var cfg = config.GetConfig()

func GetUserRepository() contractRepository.UserRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{{Entity: "UserRoles.Role"}}
	return infraRepository.NewUserRepository(cfg, preloads)
}

func GetRoleRepository() contractRepository.RoleRepository {
	return infraRepository.NewRoleRepository(cfg)
}

func GetEventRepository() contractRepository.EventRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{{Entity: "Teacher"}}
	return infraRepository.NewEventRepository(cfg, preloads)
}

func GetRegisterEventRepository() contractRepository.RegistrationRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{{Entity: "User"}, {Entity: "Event"}}
	return infraRepository.NewRegistrationRepository(cfg, preloads)
}

func GetCertificateRepository() contractRepository.CertificateRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{{Entity: "Registration"}}
	return infraRepository.NewCertificateRepository(cfg, preloads)
}
