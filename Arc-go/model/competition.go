package model

import (
	"Arc/util"
	"github.com/jinzhu/gorm"
	"time"
)

type Competition struct {
	gorm.Model
	// * 这个id不能加，加了AutoMigrate会报错
	//ID              uint             `gorm:"primary_key"`
	Title           string           `gorm:"type:varchar(100);not null"`
	Description     string           `gorm:"type:varchar(255);not null"`
	ParticipantsNum int              `gorm:"type:int;not null"`
	Category        string           `gorm:"type:varchar(100);not null"` // 比赛类型
	Tags            util.StringSlice `gorm:"type:varchar(100);not null"` // 比赛标签
	Active          int              `gorm:"type:int;not null"`          // 0:不可见 1:可见
	StartTime       time.Time        `gorm:"type:datetime;not null"`
	EndTime         time.Time        `gorm:"type:datetime;not null"`
	Password        string           `gorm:"type:varchar(100)"`
}
