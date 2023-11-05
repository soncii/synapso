package service

import (
	"github.com/labstack/echo/v4"
	"synapso/model"
	"synapso/repository"
	middles "synapso/transport/middleware"
)

func SaveRecognition(ctx echo.Context, recognition model.RecognitionDTO) (model.RecognitionDTO, error) {
	repo := repository.GetRepository()
	m := recognition.ToModel(middles.GetUserIDFromContext(ctx))
	err := repo.Save(&m).Error
	if err != nil {
		return recognition, err
	}
	data := recognition.Data
	err = repo.Transaction(func(tx *repository.Repository) error {
		for i := 0; i < len(data); i++ {
			data[i].RecognitionId = m.Id
			tx.Save(&data[i])
		}
		return nil
	})
	if err != nil {
		return model.RecognitionDTO{}, err
	}
	return m.ToDTO(data), nil
}

func GetRecognition(ctx echo.Context, id int) (model.RecognitionDTO, error) {
	repo := repository.GetRepository()
	var recognition model.Recognition
	err := repo.First(&recognition, id).Error
	if err != nil {
		return model.RecognitionDTO{}, err
	}
	var data []model.RecognitionData
	err = repo.Where("recognition_id = ?", id).Find(&data).Error
	if err != nil {
		return model.RecognitionDTO{}, err
	}
	return recognition.ToDTO(data), nil
}
