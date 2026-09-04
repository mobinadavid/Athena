package seeders

import (
	"athena/src/database"
	"athena/src/models"
	"log"
)

func SeedAdmins() {
	isActive := true
	admins := []*models.AdminModel{
		{
			IsActive:     &isActive,
			FirstName:    "Mobina",
			LastName:     "Davoudi",
			Username:     "mobina",
			Password:     []byte("!Mobina123"),
			Mobile:       "+989388835773",
			TwoFaEnabled: false,
		},
	}

	for _, admin := range admins {
		database.GetInstance().GetClient().FirstOrCreate(&admin, models.AdminModel{Mobile: admin.Mobile, Email: admin.Email})

		var superAdminRole models.RoleModel
		database.GetInstance().GetClient().Where("name = ?", models.SuperAdminRole).Find(&superAdminRole)
		database.GetInstance().GetClient().Model(&admin).Association("Roles").Replace(&superAdminRole)
	}

	log.Println("Admin Seeder executed successfully.")
}
