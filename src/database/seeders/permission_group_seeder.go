package seeders

import (
	"athena/src/database"
	"athena/src/models"
	"log"

	"gorm.io/datatypes"
)

func SeedPermissionGroup() {
	db := database.GetInstance().GetClient()

	var allPermissions []*models.PermissionModel
	if err := db.Find(&allPermissions).Error; err != nil {
		log.Fatalf("Failed to fetch permissions: %v", err)
	}
	var defaultPermissions []*models.PermissionModel
	if err := db.Where("name = ?", "admin-show").Find(&defaultPermissions).Error; err != nil {
		log.Fatalf("Failed to fetch default permissions: %v", err)
	}
	groups := []models.PermissionGroupModel{{
		Name:        "all-permissions",
		Title:       datatypes.JSON(`{"fa": "تمام دسترسی ها"}`),
		Permissions: allPermissions,
	},
		{
			Name:        "default-permission-group",
			Title:       datatypes.JSON(`{"fa": "گروه دسترسی پایه"}`),
			Permissions: defaultPermissions,
		},
	}
	for _, group := range groups {
		if err := db.FirstOrCreate(&group, models.PermissionGroupModel{Name: group.Name}).Error; err != nil {
			log.Fatalf("Failed to create permission group: %v", err)
		}

	}

	log.Println("Permission group seeded successfully.")
}
