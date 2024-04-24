package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
	Role      int    `gorm:"type:int;default:1"` // 0 超级管理员 1 管理员 2 普通用户
	TeamId    int    `gorm:"type:int;default:0"` // 绑定队伍id
}
