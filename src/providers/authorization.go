package providers

import (
	"athena/src/api/http/controllers/admins"
	"athena/src/api/http/middlewares"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"
)

func ProvideRoleRepository(db *database.Database) *repositories.RoleRepository {
	return &repositories.RoleRepository{
		DatabaseHandler: db,
	}
}

func ProvidePermissionRepository(db *database.Database) *repositories.PermissionRepository {
	return &repositories.PermissionRepository{
		DatabaseHandler: db,
	}
}

func ProvidePermissionGroupRepository(db *database.Database) *repositories.PermissionGroupRepository {
	return &repositories.PermissionGroupRepository{
		DatabaseHandler: db,
	}
}

func ProvideAuthorizationService(roleRepository *repositories.RoleRepository,
	permissionRepository *repositories.PermissionRepository,
	userRepository *repositories.UserRepository, adminRepository *repositories.AdminRepository, permissionGroup *repositories.PermissionGroupRepository) *services.AuthorizationService {
	return &services.AuthorizationService{
		RoleRepository:            roleRepository,
		PermissionRepository:      permissionRepository,
		UserRepository:            userRepository,
		AdminRepository:           adminRepository,
		PermissionGroupRepository: permissionGroup,
	}
}

func ProvideAuthorizationMiddleware(authorizationService *services.AuthorizationService) *middlewares.AuthorizationMiddleware {
	return &middlewares.AuthorizationMiddleware{
		AuthorizationService: authorizationService,
	}
}
func ProvideAdminAuthorizationController(authorizationService *services.AuthorizationService) *admins.AuthorizationController {
	return &admins.AuthorizationController{
		AuthorizationService: authorizationService,
	}
}
