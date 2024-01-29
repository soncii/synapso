package controller

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"synapso/model"
	"synapso/service"
)

func SaveRecognitionExperiment(ctx echo.Context) (err error) {
	var recognition model.RecognitionDTO
	err = ctx.Bind(&recognition)
	if recognition.Name == "" || recognition.Data == nil {
		return ctx.JSON(400, "name and data are required")
	}
	if err != nil {
		return ctx.NoContent(400)
	}
	recognition, err = service.SaveRecognition(ctx, recognition)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.JSON(201, recognition)
}

func GetRecognitionExperiment(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(400, err)
	}
	recognition, err := service.GetRecognition(ctx, intId)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.JSON(200, recognition)
}
