package controller

import (
	"net/http"
	"synapso/logger"

	"github.com/labstack/echo/v4"
)

// APIError represents
type APIError struct {
	Code    int
	Message string
}

// JSONErrorHandler is customize error handler
func JSONErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := err.Error()
	logger.GetEchoLogger().Error(err)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = err.Error()
	}

	var apierr APIError
	apierr.Code = code
	apierr.Message = msg

	if !c.Response().Committed {
		if err := c.JSON(code, apierr); err != nil {
			c.Logger().Error(err)
		}
	}
	c.Logger().Debug(err)
}
