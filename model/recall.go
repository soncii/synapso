package model

type RecallDTO struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	UserId           int      `json:"userId"`
	Stimulus         Stimulus `json:"stimulus"`
	IsSequenceMatter bool     `json:"isSequenceMatter"`
	IsFreeRecall     bool     `json:"isFreeRecall"`
}

func (r RecallDTO) ToModel(userID int) Recall {
	var recall Recall
	recall.ID = r.ID
	recall.UserID = userID
	recall.IsFreeRecall = r.IsFreeRecall
	recall.IsSequenceMatter = r.IsSequenceMatter
	recall.Type = r.Stimulus.Type
	recall.Stimuli = r.Stimulus.Stimuli
	return recall
}

func (r Recall) ToDTO() RecallDTO {
	var recallDTO RecallDTO
	recallDTO.ID = r.ID
	recallDTO.Name = r.Name
	recallDTO.UserId = r.UserID
	recallDTO.IsFreeRecall = r.IsFreeRecall
	recallDTO.IsSequenceMatter = r.IsSequenceMatter
	recallDTO.Stimulus.Type = r.Type
	recallDTO.Stimulus.Stimuli = r.Stimuli
	return recallDTO
}

type Recall struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	UserID           int       `json:"userId"`
	Type             string    `json:"type"`
	Stimuli          []Stimuli `json:"stimulus" gorm:"-"`
	IsSequenceMatter bool      `json:"isSequenceMatter"`
	IsFreeRecall     bool      `json:"isFreeRecall"`
}

type Stimulus struct {
	Type    string    `json:"type"`
	Stimuli []Stimuli `json:"stimuli"`
}

type Stimuli struct {
	ID       int    `json:"id"`
	RecallID int    `json:"recallId"`
	Data     string `json:"data"`
	Delay    int    `json:"delay"`
	Cue      string `json:"cue"`
}

type RecallResult struct {
	ID             int    `json:"id"`
	RecallID       int    `json:"recallId"`
	UserId         int    `json:"userId"`
	Response       string `json:"response"`
	TimeToComplete int    `json:"timeToComplete"`
}
