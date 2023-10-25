package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"synapso/logger"
	"synapso/model"
	"synapso/repository"
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

func SignUpUser(ctx context.Context, user model.User) (*model.User, error) {
	rep := repository.GetRepository()
	result, err := user.Create(rep)
	if err != nil {
		logger.GetEchoLogger().Error(err)
		return nil, err
	}

	return result, nil
}
