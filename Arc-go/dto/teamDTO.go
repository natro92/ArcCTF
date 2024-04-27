package dto

import (
	"Arc/model"
	"Arc/util"
)

type TeamDto struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	MemberNum     int              `json:"member_num"`
	Leader        int              `json:"leader"`
	CompetitionId util.StringSlice `json:"competition_id"`
}

func ToTeamDto(team model.Team) TeamDto {
	return TeamDto{
		ID:            team.ID,
		Name:          team.Name,
		MemberNum:     team.MemberNum,
		Leader:        team.Leader,
		CompetitionId: team.CompetitionId,
	}
}
