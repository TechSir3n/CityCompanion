package app

import (
	"github.com/TechSir3n/CityCompanion/api"
	"github.com/TechSir3n/CityCompanion/database"
)

func StartTGBot() {
	database.ConnectDB()
	api.TGBot()

	defer database.DB.Close()
}
