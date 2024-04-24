package dto

import (
	"Arc/model"
)

type ChallengeDto struct {
	ID            uint     `json:"id"`
	CompetitionID uint     `json:"competition_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Score         int      `json:"score"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	Hints         []string `json:"hints"`
}

func ToChallengeDto(challenge model.Challenge) ChallengeDto {
	return ChallengeDto{
		ID:            challenge.ID,
		CompetitionID: challenge.CompetitionID,
		Title:         challenge.Title,
		Description:   challenge.Description,
		Score:         challenge.Score,
		Category:      challenge.Category,
		Tags:          challenge.Tags,
		Hints:         challenge.Hints,
	}
}

func ToChallengesListDto(challenges []model.Challenge) []ChallengeDto {
	challengesDto := make([]ChallengeDto, 0)
	for _, challenge := range challenges {
		challengesDto = append(challengesDto, ToChallengeDto(challenge))
	}
	return challengesDto
}
