package migration

import (
	"github.com/ybkuroki/go-webapp-project-template/config"
	"github.com/ybkuroki/go-webapp-project-template/model"
	"github.com/ybkuroki/go-webapp-project-template/repository"
)

// InitMasterData creates the master data used in this application.
func InitMasterData(config *config.Config) {
	if config.Extension.MasterGenerator {
		rep := repository.GetRepository()

		r := model.User{Email: "admin@gmail.com", Role: enums.}
		_, _ = r.Create(rep)
		a := model.NewAccountWithPlainPassword("test", "test", r)
		_, _ = a.Create(rep)
	}
}
