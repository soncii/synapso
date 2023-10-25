package migration

import (
	"synapso/config"
	"synapso/enums"
	"synapso/model"
	"synapso/repository"
)

// InitMasterData creates the master data used in this application.
func InitMasterData(config *config.Config) {
	if config.Extension.MasterGenerator {
		rep := repository.GetRepository()

		user := model.User{Email: "admin@gmail.com", Role: enums.RESEARCHER, Password: string(model.HashPassword("admin"))}
		_, _ = user.Create(rep)
		subject := model.User{Email: "subject@gmail.com", Role: enums.SUBJECT, Password: string(model.HashPassword("subject"))}
		_, _ = subject.Create(rep)
	}
}
