package dto

import (
	"Arc/model"
	"Arc/util"
	"time"
)

type CompetitionDto struct {
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	ParticipantsNum int              `json:"participants_num"`
	Category        string           `json:"category"`
	Tags            util.StringSlice `json:"tags"`
	StartTime       time.Time        `json:"start_time"`
	EndTime         time.Time        `json:"end_time"`
}

func ToCompetitionDto(competition model.Competition) CompetitionDto {
	return CompetitionDto{
		Title:           competition.Title,
		Description:     competition.Description,
		ParticipantsNum: competition.ParticipantsNum,
		Category:        competition.Category,
		Tags:            competition.Tags,
		StartTime:       competition.StartTime,
		EndTime:         competition.EndTime,
	}
}

func ToCompetitionsListDto(competitions []model.Competition) []CompetitionDto {
	competitionsDto := make([]CompetitionDto, 0)
	for _, competition := range competitions {
		competitionsDto = append(competitionsDto, ToCompetitionDto(competition))
	}
	return competitionsDto
}
