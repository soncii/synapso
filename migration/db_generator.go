package migration

import (
	config "synapso/config"
	"synapso/model"
	"synapso/repository"
)

// CreateDatabase creates the tables used in this application.
func CreateDatabase(config *config.Config) {
	if config.Database.Migration {
		db := repository.GetDB()
		db.DropTable(&model.User{})
		db.AutoMigrate(&model.User{})
	}
}
