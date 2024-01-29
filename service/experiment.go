package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
	"synapso/logger"
	"synapso/model"
	"synapso/repository"
	middles "synapso/transport/middleware"
)

func ListExperiments(ctx echo.Context) (result model.ExperimentList, err error) {
	repo := repository.GetRepository()
	var recalls []model.Recall
	err = repo.Find(&recalls).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}
	for i := 0; i < len(recalls); i++ {
		repo.Find(&recalls[i].Stimuli, "recall_id = ?", recalls[i].ID)
		result.Recall = append(result.Recall, recalls[i].ToDTO())
	}

	var recognition []model.Recognition
	err = repo.Find(&recognition).Error
	for i := 0; i < len(recognition); i++ {
		var data []model.RecognitionData
		repo.Find(&data, "recognition_id = ?", recognition[i].Id)
		result.Recognition = append(result.Recognition, recognition[i].ToDTO(data))
	}
	return
}

func SaveExperimentResult(ctx echo.Context, recognitionResult model.ExperimentResultDTO) error {
	repo := repository.GetRepository()
	recognitionResult.UserId = middles.GetUserIDFromContext(ctx)
	var user model.User
	repo.Where("id=?", recognitionResult.UserId).First(&user)
	recognitionResult.Name = user.FirstName + " " + user.Surname
	m := recognitionResult.ToModel()
	return repo.Save(&m).Error
}

func GetExperimentResult(ctx echo.Context, id int) (model.ExperimentResultDTO, error) {
	repo := repository.GetRepository()
	var recognitionResult model.ExperimentResult
	err := repo.First(&recognitionResult, id).Error

	return recognitionResult.ToDTO(), err
}

func GetExperimentResultByExperimentIdAndType(ctx echo.Context, experimentId int, experimentType string) (body bytes.Buffer, err error) {
	var results []model.ExperimentResult
	repo := repository.GetRepository()
	err = repo.Where("recognition_id = ? AND type = ?", experimentId, experimentType).Find(&results).Error
	var recognition model.RecognitionDTO
	if experimentType == "recognition" {
		recognition, err = GetRecognition(ctx, experimentId)
		if err != nil {
			return body, err
		}
	}
	logger.GetEchoLogger().Info("results: ", results, "recognition: ", recognition)
	var w bytes.Buffer
	writer := csv.NewWriter(&w)

	// Write the CSV header
	header := []string{"Id", "Subject Name", "TimeToComplete", "%correct"}
	for _, data := range recognition.Data {
		header = append(header, data.Displayed+"/"+data.Hidden)
	}
	err = writer.Write(header)
	if err != nil {
		logger.GetEchoLogger().Error(err)
	}

	for _, result := range results {
		// Convert Response to a comma-separated string
		record := []string{
			strconv.Itoa(result.Id),
			result.Name,
			strconv.Itoa(result.TimeToComplete),
		}
		var correct int
		var words []string
		for i, data := range strings.Split(result.Response, ",") {
			words = append(words, data)
			if data == recognition.Data[i].Displayed {
				correct += 1
			}
		}
		if len(result.Response) == 0 {
			record = append(record, "0.00")
		} else {
			percentageCompleted := float64(correct*100) / float64(len(strings.Split(result.Response, ",")))
			formattedPercentage := fmt.Sprintf("%.2f", percentageCompleted)
			record = append(record, formattedPercentage)
		}
		record = append(record, words...)

		err := writer.Write(record)
		if err != nil {
			logger.GetEchoLogger().Error(err)
		}
	}
	writer.Flush()
	logger.GetEchoLogger().Info("w: ", w.String())
	return w, nil
}

func GetExperimentResultsByUserID(ctx echo.Context, userID int) ([]model.ExperimentResultDTO, error) {
	repo := repository.GetRepository()
	var recognitionResults []model.ExperimentResult
	err := repo.Where("user_id = ?", userID).Find(&recognitionResults).Error
	if err != nil {
		return nil, err
	}
	var recognitionResultsDTO []model.ExperimentResultDTO
	for _, recognitionResult := range recognitionResults {
		recognitionResultsDTO = append(recognitionResultsDTO, recognitionResult.ToDTO())
	}
	return recognitionResultsDTO, nil
}
