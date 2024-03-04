package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"sort"
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

func GetUniformExperimentList(ctx echo.Context) (result []model.ExperimentCommon, err error) {
	repo := repository.GetRepository()
	var recalls []model.Recall
	err = repo.Find(&recalls).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}
	for i := 0; i < len(recalls); i++ {
		var exp model.ExperimentCommon
		exp.Id = recalls[i].ID
		exp.Name = recalls[i].Name
		exp.Type = "recall"
		exp.DistractionType = recalls[i].DistractionType
		exp.CreatedAt = model.CustomTime{Time: recalls[i].CreatedAt}
		exp.StimulusType = recalls[i].Type
		repo.Db.Model(&model.ExperimentResult{}).Where("recognition_id = ? AND type = ?", recalls[i].ID, "recall").Count(&exp.UsersResponded)
		result = append(result, exp)
	}

	var recognition []model.Recognition
	err = repo.Find(&recognition).Error
	for i := 0; i < len(recognition); i++ {
		var exp model.ExperimentCommon
		exp.Id = recognition[i].Id
		exp.Name = recognition[i].Name
		exp.Type = "recognition"
		exp.DistractionType = recognition[i].DistractionType
		exp.CreatedAt = model.CustomTime{Time: recognition[i].CreatedAt}
		exp.StimulusType = recognition[i].Type
		repo.Db.Model(&model.ExperimentResult{}).Where("recognition_id = ? AND type = ?", recognition[i].Id, "recognition").Count(&exp.UsersResponded)
		result = append(result, exp)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})
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

func GetRecognitionResultByExperimentId(ctx echo.Context, experimentId int) (body bytes.Buffer, filename string, err error) {
	var results []model.ExperimentResult
	repo := repository.GetRepository()
	err = repo.Where("recognition_id = ? AND type = ?", experimentId, "recognition").Find(&results).Error
	recognition, err := GetRecognition(ctx, experimentId)
	if err != nil {
		return
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
	return w, recognition.Name, nil
}

func GetRecallResultByExperimentId(ctx echo.Context, experimentId int) (body bytes.Buffer, filename string, err error) {
	var results []model.ExperimentResult
	repo := repository.GetRepository()
	err = repo.Where("recognition_id = ? AND type = ?", experimentId, "recall").Find(&results).Error
	recall, err := GetRecallExperiment(context.TODO(), experimentId)
	if err != nil {
		return
	}
	logger.GetEchoLogger().Info("results: ", results, "recall: ", recall)
	var w bytes.Buffer
	writer := csv.NewWriter(&w)

	// Write the CSV header
	header := []string{"Id", "Subject Name", "TimeToComplete", "%correct"}
	for _, data := range recall.Stimulus.Stimuli {
		header = append(header, data.Data)
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
		for _, data := range strings.Split(result.Response, ",") {
			words = append(words, data)
			for _, cor := range recall.Stimulus.Stimuli {
				if data == cor.Data {
					correct += 1
					break
				}
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
	return w, recall.Name, nil
}

func GetExperimentResultsByUserID(ctx echo.Context, userID int) ([]model.ExperimentResultDTO, error) {

	repo := repository.GetRepository()
	var recognitionResults = make([]model.ExperimentResult, 0)
	err := repo.Where("user_id = ?", userID).Find(&recognitionResults).Error
	if err != nil {
		return nil, err
	}
	var recognitionResultsDTO = make([]model.ExperimentResultDTO, 0)
	for _, recognitionResult := range recognitionResults {
		recognitionResultsDTO = append(recognitionResultsDTO, recognitionResult.ToDTO())
	}
	return recognitionResultsDTO, nil
}

func DeleteAllExperiment(ctx echo.Context) error {
	repo := repository.GetRepository()
	repo.Delete(model.Recognition{})
	repo.Delete(model.RecognitionData{})
	repo.Delete(model.Recall{})
	repo.Delete(model.Stimulus{})
	repo.Delete(model.ExperimentResult{})
	return nil
}
