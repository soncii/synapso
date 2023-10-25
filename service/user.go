package service

import (
	"github.com/ybkuroki/go-webapp-project-template/logger"
	"github.com/ybkuroki/go-webapp-project-template/model"
	"github.com/ybkuroki/go-webapp-project-template/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateByUsernameAndPassword authenticates by using username and plain text password.
func AuthenticateByUsernameAndPassword(username string, password string) (bool, *model.User) {
	rep := repository.GetRepository()
	User := model.User{}
	result, err := User.FindByEmail(rep, username)
	if err != nil {
		logger.GetEchoLogger().Error(err)
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		logger.GetEchoLogger().Error(err)
		return false, nil
	}

	return true, result
}
