package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"synapso/model"
	"synapso/repository"
)

func SaveRecallExperiment(ctx context.Context, recall model.Recall) (int, error) {
	repo := repository.GetRepository()
	err := repo.Save(&recall).Error
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(recall.Stimuli); i++ {
		recall.Stimuli[i].RecallID = recall.ID
	}
	fmt.Println("saving stimuli")
	fmt.Println(recall.Stimuli)

	err = repo.Db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(recall.Stimuli); i++ {
			err := tx.Save(&recall.Stimuli[i]).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return recall.ID, err
}

func GetRecallExperiment(ctx context.Context, id int) (model.RecallDTO, error) {
	repo := repository.GetRepository()
	var recallExperiment model.Recall
	err := repo.First(&recallExperiment, id).Error
	fmt.Println("fetched main recall experiment")
	var st []model.Stimuli
	repo.Db.Where("recall_id = ?", recallExperiment.ID).Find(&st)
	fmt.Println("fetched stimuli")
	fmt.Println(st)
	recallExperiment.Stimuli = st
	return recallExperiment.ToDTO(), err
}
