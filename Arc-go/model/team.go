package model

import (
	"Arc/util"
	"github.com/jinzhu/gorm"
)

type Team struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null"`
	//Members    map[string]*User `gorm:"type:varchar(100);not null"` 似乎不用加 ，在user里面加一个teamid
	Leader        int              `gorm:"type:int;not null"` // 绑定队长id
	MemberNum     int              `gorm:"type:int;not null"` // 队伍人数
	InviteCode    string           `gorm:"type:varchar(100)"` // 邀请码
	CompetitionId util.StringSlice `gorm:"type:varchar(100)"` // 比赛id
}
