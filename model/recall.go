package model

import "time"

type RecallDTO struct {
	ID                   int      `json:"id"`
	Name                 string   `json:"name"`
	UserId               int      `json:"userId"`
	InstructionText      string   `json:"instructionText"`
	IsDistractionEnabled bool     `json:"isDistractionEnabled"`
	DistractionType      string   `json:"distractionType"`
	DistractionText      string   `json:"distractionText"`
	DistractionDuration  int      `json:"distractionDuration"`
	InterStimuliDelay    int      `json:"interStimuliDelay"`
	IsFreeRecall         bool     `json:"isFreeRecall"`
	Stimulus             Stimulus `json:"stimulus"`
}

func (r RecallDTO) ToModel(userID int) Recall {
	var recall Recall
	recall.ID = r.ID
	recall.Name = r.Name
	recall.UserID = userID
	recall.CreatedAt = time.Now().UTC()
	recall.InstructionText = r.InstructionText
	recall.IsFreeRecall = r.IsFreeRecall
	recall.Type = r.Stimulus.Type
	recall.DistractionDuration = r.DistractionDuration
	recall.DistractionText = r.DistractionText
	recall.DistractionType = r.DistractionType
	recall.InterStimuliDelay = r.InterStimuliDelay
	recall.IsDistractionEnabled = r.IsDistractionEnabled
	recall.Stimuli = r.Stimulus.Stimuli
	return recall
}

func (r Recall) ToDTO() RecallDTO {
	var recallDTO RecallDTO
	recallDTO.ID = r.ID
	recallDTO.Name = r.Name
	recallDTO.InstructionText = r.InstructionText
	recallDTO.UserId = r.UserID
	recallDTO.IsFreeRecall = r.IsFreeRecall
	recallDTO.IsDistractionEnabled = r.IsDistractionEnabled
	recallDTO.DistractionType = r.DistractionType
	recallDTO.DistractionText = r.DistractionText
	recallDTO.DistractionDuration = r.DistractionDuration
	recallDTO.InterStimuliDelay = r.InterStimuliDelay
	recallDTO.Stimulus.Type = r.Type
	recallDTO.Stimulus.Stimuli = r.Stimuli
	return recallDTO
}

type Recall struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	UserID               int       `json:"userId"`
	CreatedAt            time.Time `json:"createdAt"`
	Type                 string    `json:"type"`
	InstructionText      string    `json:"instructionText"`
	Stimuli              []Stimuli `json:"stimulus" gorm:"-"`
	IsDistractionEnabled bool      `json:"isDistractionEnabled"`
	DistractionType      string    `json:"distractionType"`
	DistractionText      string    `json:"distractionText"`
	DistractionDuration  int       `json:"distractionDuration"`
	InterStimuliDelay    int       `json:"interStimuliDelay"`
	IsFreeRecall         bool      `json:"isFreeRecall"`
}

type Stimulus struct {
	Type    string    `json:"type"`
	Stimuli []Stimuli `json:"stimuli"`
}

type Stimuli struct {
	ID       int     `json:"id"`
	RecallID int     `json:"recallId"`
	Data     string  `json:"data"`
	Delay    int     `json:"delay"`
	Cue      *string `json:"cue"`
}

type RecallResult struct {
	ID             int    `json:"id"`
	RecallID       int    `json:"recallId"`
	UserId         int    `json:"userId"`
	Response       string `json:"response"`
	TimeToComplete int    `json:"timeToComplete"`
}
