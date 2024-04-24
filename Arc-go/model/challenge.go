package model

import (
	"Arc/util"
	"github.com/jinzhu/gorm"
	"time"
)

// * 定义 Challenge 题目
type Challenge struct {
	gorm.Model
	CompetitionID   uint             `gorm:"type:int;not null"`
	Title           string           `gorm:"type:varchar(100);not null"`
	Description     string           `gorm:"type:varchar(255)"`
	Score           int              `gorm:"type:int;not null"`
	MinScore        int              `gorm:"type:int;not null"`
	MaxScore        int              `gorm:"type:int;not null"`
	ContainerMirror string           `gorm:"type:varchar(100)"`
	Flag            string           `gorm:"type:varchar(100);not null"`
	Attachment      util.StringSlice `gorm:"type:varchar(100)"`
	Category        string           `gorm:"type:varchar(100);not null"`
	Tags            util.StringSlice `gorm:"type:varchar(100);not null"`
	Hints           util.StringSlice `gorm:"type:varchar(255);not null"`
	Visible         int              `gorm:"type:int;not null"` // 0:不可见 1:可见
}

// * 定义用户 Submission 提交
type Submission struct {
	gorm.Model
	ID          uint      `gorm:"primary_key"`
	UserID      uint      `gorm:"type:int;not null"`
	TeamID      uint      `gorm:"type:int;not null"`
	Answer      string    `gorm:"type:varchar(100);not null"`
	SubmitTime  time.Time `gorm:"type:datetime;not null"`
	GameID      uint      `gorm:"type:int;not null"`
	GameEndTime time.Time `gorm:"type:datetime;not null"`
	ChallengeID uint      `gorm:"type:int;not null"`
	Status      int       `gorm:"type:int;not null"` // 0:未解决 1:已解决
	Flag        string    `gorm:"type:varchar(100);not null"`
}

// * 定义用户 Solve 正确的解题
type Solve struct {
	gorm.Model
	ID          uint `gorm:"primary_key"`
	UserID      uint `gorm:"type:int;not null"`
	ChallengeID uint `gorm:"type:int;not null"`
	Time        int  `gorm:"type:int;not null"` // 解题时间
}
