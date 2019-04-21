package main

import (
	"github.com/jinzhu/gorm"
)

// type UserInfo struct {
// 	Username string `json:"username" required`
// 	Password string `json:"password" required`
// }

// type Response struct {
// 	token string `json:"token"`
// }

type UserInfo struct {
	gorm.Model
	Username string   `gorm:"type:varchar(16);not null;" json:"username required"`
	Password string   `gorm:"type:varchar(128);not null;" json:"password required"`
	RoleID   uint     `gorm:"not null;" json:"role_id required"`
	Role     RoleInfo `gorm:"foreignkey:RoleID"`
}

type RoleInfo struct {
	gorm.Model
	Name string `gorm:"type:varchar(16);not null;" json:"role_name required"`
}

type ProjectInfo struct {
	gorm.Model
	Name string `gorm:"type:varchar(32);not null;" json:"project_name required"`
}

type ProjectEnv struct {
	gorm.Model
	ProjectID uint        `gorm:"type:varchar(32);not null;" json:"project_id required"`
	Project   ProjectInfo `gorm:"foreignkey:ProjectID"`
	EnvKey    string      `gorm:"type:varchar(32);not null;" json:"env_key required"`
	EnvValue  string      `gorm:"type:varchar(128);not null;" json:"env_value required"`
	IsSecret  bool        `gorm:"not null;" json:"is_secret required"`
}
