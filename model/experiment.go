package model

import "strings"

type ExperimentList struct {
	Recall      []RecallDTO      `json:"recall"`
	Recognition []RecognitionDTO `json:"recognition"`
}

type ExperimentResult struct {
	Id             int    `json:"id"`
	Type           string `json:"experimentType"`
	RecognitionId  int    `json:"experimentId"`
	UserId         int    `json:"userId"`
	Response       string `json:"response"`
	TimeToComplete int    `json:"timeToComplete"`
}

type ExperimentResultDTO struct {
	Id             int      `json:"id"`
	Type           string   `json:"type"`
	RecognitionId  int      `json:"recognitionId"`
	UserId         int      `json:"userId"`
	Response       []string `json:"response"`
	TimeToComplete int      `json:"timeToComplete"`
}

func (d ExperimentResultDTO) ToModel() (result ExperimentResult) {
	result.Id = d.Id
	result.Type = d.Type
	result.RecognitionId = d.RecognitionId
	result.UserId = d.UserId
	result.TimeToComplete = d.TimeToComplete
	result.Response = strings.Join(d.Response, ",")
	return
}

func (d ExperimentResult) ToDTO() (result ExperimentResultDTO) {
	result.Id = d.Id
	result.Type = d.Type
	result.RecognitionId = d.RecognitionId
	result.UserId = d.UserId
	result.TimeToComplete = d.TimeToComplete
	result.Response = strings.Split(d.Response, ", ")
	return
}
