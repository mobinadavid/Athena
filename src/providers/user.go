package providers

import (
	"athena/src/api/http/controllers/admins"
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

func ProvideAdminUserController(userService *services.UserService) *admins.UserController {
	return &admins.UserController{
		UserService: userService,
	}
}
