package providers

import (
	"athena/src/api/http/controllers/admins"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"
)

func ProvideAdminRepository(db *database.Database) *repositories.AdminRepository {
	return &repositories.AdminRepository{
		DatabaseHandler: db,
	}
}

func ProvideAdminService(adminRepository *repositories.AdminRepository, roleRepository *repositories.RoleRepository, OTPService *services.OTPService) *services.AdminService {
	return &services.AdminService{
		AdminRepository: adminRepository,
		RoleRepository:  roleRepository,
		OTPService:      OTPService,
	}
}

func ProvideAdminController(adminService *services.AdminService) *admins.AdminController {
	return &admins.AdminController{
		AdminService: adminService,
	}
}
