package controller

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"synapso/logger"
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

func GetUniformExperimentList(ctx echo.Context) (err error) {
	experiments, err := service.GetUniformExperimentList(ctx)
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
	logger.GetEchoLogger().Info("experimentType: "+experimentType, "id: "+id)
	var body bytes.Buffer
	var filename string
	if experimentType == "recognition" {
		body, filename, err = service.GetRecognitionResultByExperimentId(ctx, intId)
		if err != nil {
			return ctx.JSON(500, err)
		}
	} else {
		body, filename, err = service.GetRecallResultByExperimentId(ctx, intId)
		if err != nil {
			return ctx.JSON(500, err)
		}
	}

	if err != nil {
		return ctx.JSON(500, err)
	}
	ctx.Response().Header().Set("Content-Type", "text/csv")
	ctx.Response().Header().Set("Content-Description", "File Transfer")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+filename+"_results.csv")
	ctx.Response().Header().Set("Content-Description", "File Transfer")
	return ctx.Blob(http.StatusOK, "text/csv", body.Bytes())

}

func GetExperimentResults(ctx echo.Context) (err error) {
	userId := middles.GetUserIDFromContext(ctx)
	recognitionResult, err := service.GetExperimentResultsByUserID(ctx, userId)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.JSON(200, recognitionResult)
}

func DeleteAllExperiment(ctx echo.Context) (err error) {
	err = service.DeleteAllExperiment(ctx)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.NoContent(200)
	return
}
