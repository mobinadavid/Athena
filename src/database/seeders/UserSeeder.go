package seeders

import (
	"athena/src/database"
	"athena/src/models"
	"log"
)

func SeedUsers() {
	users := []*models.UserModel{
		{
			FirstName:            "مبینا",
			LastName:             "داودی",
			NationalIdentityCode: "0441226086",
			Mobile:               "09388835773",
			Password:             []byte("!Mobina123"),
			Email:                "h@gmail.com",
		},
	}

	for _, user := range users {
		database.GetInstance().GetClient().FirstOrCreate(&user, models.UserModel{Mobile: user.Mobile, NationalIdentityCode: user.NationalIdentityCode})
	}

	log.Println("User Seeder executed successfully.")
}
