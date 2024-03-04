package model

import (
	"fmt"
	"strings"
	"time"
)

const layout = "02.01.2006 15:04"

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(layout))), nil
}
func (ct CustomTime) UnmarshalJSON(b []byte) error {
	// Remove the quotes from the JSON string
	s := string(b)
	s = s[1 : len(s)-1] // slice off the quotes

	// Parse the time string based on the custom format
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	ct.Time = t
	return nil
}

type ExperimentList struct {
	Recall      []RecallDTO      `json:"recall"`
	Recognition []RecognitionDTO `json:"recognition"`
}

type ExperimentCommon struct {
	Id              int        `json:"id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	DistractionType string     `json:"distractionType"`
	CreatedAt       CustomTime `json:"createdAt"`
	StimulusType    string     `json:"stimulusType"`
	UsersResponded  int        `json:"usersResponded"`
}

type ExperimentResult struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"experimentType"`
	RecognitionId  int    `json:"experimentId"`
	UserId         int    `json:"userId"`
	Response       string `json:"response"`
	TimeToComplete int    `json:"timeToComplete"`
}

type ExperimentResultDTO struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	RecognitionId  int      `json:"experimentId"`
	UserId         int      `json:"userId"`
	Response       []string `json:"response"`
	TimeToComplete int      `json:"timeToComplete"`
}

func (d ExperimentResultDTO) ToModel() (result ExperimentResult) {
	result.Id = d.Id
	result.Name = d.Name
	result.Type = d.Type
	result.RecognitionId = d.RecognitionId
	result.UserId = d.UserId
	result.TimeToComplete = d.TimeToComplete
	result.Response = strings.Join(d.Response, ",")
	return
}

func (d ExperimentResult) ToDTO() (result ExperimentResultDTO) {
	result.Id = d.Id
	result.Name = d.Name
	result.Type = d.Type
	result.RecognitionId = d.RecognitionId
	result.UserId = d.UserId
	result.TimeToComplete = d.TimeToComplete
	result.Response = strings.Split(d.Response, ", ")
	return
}
