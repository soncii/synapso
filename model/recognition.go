package model

type RecognitionDTO struct {
	Id     int               `json:"id"`
	Name   string            `json:"name"`
	UserId int               `json:"userId"`
	Type   string            `json:"type"`
	Data   []RecognitionData `gorm:"-" json:"data"`
}

func (r RecognitionDTO) ToModel(userID int) Recognition {
	var recognition Recognition
	recognition.Id = r.Id
	recognition.Name = r.Name
	recognition.UserId = userID
	recognition.Type = r.Type
	return recognition
}

type Recognition struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func (r Recognition) ToDTO(Data []RecognitionData) RecognitionDTO {
	var recognitionDTO RecognitionDTO
	recognitionDTO.Id = r.Id
	recognitionDTO.Name = r.Name
	recognitionDTO.UserId = r.UserId
	recognitionDTO.Type = r.Type
	recognitionDTO.Data = Data
	return recognitionDTO
}

type RecognitionData struct {
	Id            int    `json:"id"`
	RecognitionId int    `json:"recognitionId"`
	Displayed     string `json:"displayed"`
	Hidden        string `json:"hidden"`
	Delay         int    `json:"delay"`
}
