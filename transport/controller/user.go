package controller

import (
	"fmt"
	"net/http"
	"synapso/enums"
	"synapso/logger"
	"synapso/model"
	"synapso/service"
	"synapso/transport/middleware"

	"github.com/labstack/echo/v4"
)

// GetLoginAccount returns the account data of logged in user.
func GetLoginAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "dummy")
	}
}

func UserSignUp(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			logger.GetEchoLogger().Error(err)
		}
	}()
	var user model.UserCreate
	err = ctx.Bind(&user)
	if err != nil {
		return err
	}
	_, err = user.Validate()
	if err != nil {
		return err
	}
	result, err := service.SignUpUser(ctx.Request().Context(), user.ToModel())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, result)
}

func UserLogin(ctx echo.Context) error {
	var auth model.Auth
	err := ctx.Bind(&auth)
	if err != nil || auth.Email == "" || auth.Password == "" {
		ctx.NoContent(http.StatusBadRequest)
		return nil
	}

	authenticate, u := service.AuthenticateByUsernameAndPassword(auth.Email, auth.Password)
	if !authenticate {
		return ctx.NoContent(http.StatusUnauthorized)
	}
	if u.Role != auth.Role {
		return ctx.JSON(http.StatusUnauthorized, fmt.Sprintf("you are %s, but you are trying to log in into %s application", u.Role, auth.Role))
	}

	token, err := middles.GenerateToken(u)
	if err != nil {
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	ctx.Response().Header().Set("Authorization", "Bearer "+token)
	if u.Role == enums.RESEARCHER {
		return ctx.JSON(http.StatusOK, "Bearer "+token)
	}
	return ctx.JSON(http.StatusOK, u)

}
