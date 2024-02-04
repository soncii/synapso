package model

type RecognitionDTO struct {
	Id                   int               `json:"id"`
	Name                 string            `json:"name"`
	UserId               int               `json:"userId"`
	Type                 string            `json:"type"`
	IsDistractionEnabled bool              `json:"isDistractionEnabled"`
	DistractionType      string            `json:"distractionType"`
	DistractionText      string            `json:"distractionText"`
	DistractionDuration  int               `json:"distractionDuration"`
	InterStimuliDelay    int               `json:"interStimuliDelay"`
	Data                 []RecognitionData `gorm:"-" json:"data"`
}

func (r RecognitionDTO) ToModel(userID int) Recognition {
	var recognition Recognition
	recognition.Id = r.Id
	recognition.Name = r.Name
	recognition.UserId = userID
	recognition.Type = r.Type
	recognition.InterStimuliDelay = r.InterStimuliDelay
	recognition.IsDistractionEnabled = r.IsDistractionEnabled
	recognition.DistractionType = r.DistractionType
	recognition.DistractionDuration = r.DistractionDuration
	recognition.DistractionText = r.DistractionText
	return recognition
}

type Recognition struct {
	Id                   int    `json:"id"`
	UserId               int    `json:"userId"`
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	InterStimuliDelay    int    `json:"interStimuliDelay"`
	IsDistractionEnabled bool   `json:"isDistractionEnabled"`
	DistractionDuration  int    `json:"distractionDuration"`
	DistractionType      string `json:"distractionType"`
	DistractionText      string `json:"distractionText"`
}

func (r Recognition) ToDTO(Data []RecognitionData) RecognitionDTO {
	var recognitionDTO RecognitionDTO
	recognitionDTO.Id = r.Id
	recognitionDTO.Name = r.Name
	recognitionDTO.DistractionDuration = r.DistractionDuration
	recognitionDTO.UserId = r.UserId
	recognitionDTO.Type = r.Type
	recognitionDTO.Data = Data
	recognitionDTO.InterStimuliDelay = r.InterStimuliDelay
	recognitionDTO.IsDistractionEnabled = r.IsDistractionEnabled
	recognitionDTO.DistractionType = r.DistractionType
	recognitionDTO.DistractionText = r.DistractionText
	return recognitionDTO
}

type RecognitionData struct {
	Id            int    `json:"id"`
	RecognitionId int    `json:"recognitionId"`
	Displayed     string `json:"displayed"`
	Hidden        string `json:"hidden"`
	Delay         int    `json:"delay"`
}
