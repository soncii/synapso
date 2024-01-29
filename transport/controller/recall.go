package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"synapso/logger"
	"synapso/model"
	"synapso/service"
	middles "synapso/transport/middleware"
)

func SaveRecall(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			logger.GetEchoLogger().Error(err)
		}
	}()
	var recall model.RecallDTO
	err = ctx.Bind(&recall)
	if recall.Name == "" {
		return ctx.JSON(400, "name and data are required")
	}
	if err != nil {
		return err
	}
	id, err := service.SaveRecallExperiment(ctx.Request().Context(), recall.ToModel(middles.GetUserIDFromContext(ctx)))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, id)
}

func GetRecall(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			logger.GetEchoLogger().Error(err)
		}
	}()
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	recall, err := service.GetRecallExperiment(ctx.Request().Context(), intId)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, recall)
}
