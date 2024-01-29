package controller

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"synapso/model"
	"synapso/service"
	middles "synapso/transport/middleware"
)

func ListExperiments(ctx echo.Context) (err error) {
	experiments, err := service.ListExperiments(ctx)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.JSON(200, experiments)
}

func SaveExperimentResult(ctx echo.Context) (err error) {
	var recognitionResult model.ExperimentResultDTO
	err = ctx.Bind(&recognitionResult)
	if err != nil {
		return ctx.JSON(400, err)
	}
	err = service.SaveExperimentResult(ctx, recognitionResult)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.NoContent(201)
}

func GetExperimentResult(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(400, err)
	}
	recognitionResult, err := service.GetExperimentResult(ctx, intId)
	if err != nil {
		return ctx.JSON(500, err)
	}

	return ctx.JSON(200, recognitionResult)
}

func GetExperimentResultCsv(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	experimentType := ctx.QueryParam("experimentType")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(400, err)
	}
	recognitionResult, err := service.GetExperimentResultByExperimentIdAndType(ctx, intId, experimentType)
	if err != nil {
		return ctx.JSON(500, err)
	}
	ctx.Response().Header().Set("Content-Type", "text/csv")
	ctx.Response().Header().Set("Content-Description", "File Transfer")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=experimentResult_"+id+".csv")
	ctx.Response().Write(recognitionResult.Bytes())
	return nil
}

func GetExperimentResults(ctx echo.Context) (err error) {
	userId := middles.GetUserIDFromContext(ctx)
	recognitionResult, err := service.GetExperimentResultsByUserID(ctx, userId)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.JSON(200, recognitionResult)
}
