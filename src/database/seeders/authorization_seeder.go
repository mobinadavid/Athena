package seeders

import (
	"athena/src/database"
	"athena/src/models"
	"log"
)

func SeedAuthorization() {
	db := database.GetInstance().GetClient()

	roles := []struct {
		Title            string
		Name             string
		PermissionGroups []string // list of permission group names
	}{
		{
			Title:            "دسترسی پایه",
			Name:             models.DefaultRole,
			PermissionGroups: []string{"default-permission-group"}, // assign by group name
		},
		{
			Title:            "سوپر ادمین",
			Name:             models.SuperAdminRole,
			PermissionGroups: []string{"all-permissions"}, // assign by group name
		},
	}

	for _, roleData := range roles {
		var role models.RoleModel

		if err := db.
			Where("name = ?", roleData.Name).
			Assign(models.RoleModel{
				Title: roleData.Title,
			}).
			FirstOrCreate(&role, models.RoleModel{
				Name:  roleData.Name,
				Title: roleData.Title,
			}).Error; err != nil {
			log.Fatalf("Failed to seed role %s: %v", roleData.Name, err)
		}

		// Find permission groups
		var groups []models.PermissionGroupModel
		if err := db.Where("name IN ?", roleData.PermissionGroups).Find(&groups).Error; err != nil {
			log.Fatalf("Failed to fetch permission groups: %v", err)
		}

		// Assign groups to role
		if err := db.Model(&role).Association("PermissionGroups").Replace(&groups); err != nil {
			log.Fatalf("Failed to assign permission groups to role %s: %v", roleData.Name, err)
		}
	}

	log.Println("Authorization Seeder executed successfully.")
}
