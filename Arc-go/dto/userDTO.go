package dto

import (
	"Arc/model"
)

type UserDto struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Role      int    `json:"role"`
	TeamId    int    `json:"team_id"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		ID:        user.ID,
		Name:      user.Name,
		Telephone: user.Telephone,
		Role:      user.Role,
		TeamId:    user.TeamId,
	}
}
