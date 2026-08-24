package providers

import (
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"
)

func ProvideUserRepository(db *database.Database) *repositories.UserRepository {
	return &repositories.UserRepository{
		DatabaseHandler: db,
	}
}

func ProvideUserService(userRepository *repositories.UserRepository) *services.UserService {
	return &services.UserService{
		UserRepository: userRepository,
	}
}
