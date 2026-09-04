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

	groups := []models.PermissionGroupModel{
		{
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
		var existing models.PermissionGroupModel
		if err := db.Where("name = ?", group.Name).
			Assign(models.PermissionGroupModel{Title: group.Title}).
			FirstOrCreate(&existing, models.PermissionGroupModel{Name: group.Name, Title: group.Title}).Error; err != nil {
			log.Fatalf("Failed to create permission group: %v", err)
		}

		if err := db.Model(&existing).Association("Permissions").Replace(group.Permissions); err != nil {
			log.Fatalf("Failed to assign permissions to group %s: %v", group.Name, err)
		}
	}

	log.Println("Permission group seeded successfully.")
}
