package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
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
	m := recognitionResult.ToModel()
	return repo.Save(&m).Error
}

func GetExperimentResult(ctx echo.Context, id int) (model.ExperimentResultDTO, error) {
	repo := repository.GetRepository()
	var recognitionResult model.ExperimentResult
	err := repo.First(&recognitionResult, id).Error

	return recognitionResult.ToDTO(), err
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
