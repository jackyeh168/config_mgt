package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserInfo struct {
	gorm.Model
	Username string   `gorm:"type:varchar(16);not null;unique_index" json:"username" binding:"required"`
	Password string   `gorm:"type:varchar(128);not null;" json:"password" binding:"required"`
	RoleID   uint     `sql:"type:int unsigned" gorm:"not null;" json:"role_id"`
	Role     RoleInfo `gorm:"foreignkey:RoleID"`
}

type RoleInfo struct {
	gorm.Model
	Name string `gorm:"type:varchar(16);not null;unique_index" json:"role_name"`
}

type ProjectInfo struct {
	gorm.Model
	Name string `gorm:"type:varchar(32);not null;unique_index" json:"project_name" binding:"required"`
}

type ProjectEnv struct {
	gorm.Model
	ProjectID uint        `sql:"type:int unsigned" gorm:"not null;" json:"project_id" binding:"required"`
	Project   ProjectInfo `gorm:"foreignkey:ProjectID"`
	EnvKey    string      `gorm:"type:varchar(32);not null;unique_index" json:"env_key" binding:"required"`
	EnvValue  string      `gorm:"type:varchar(128);not null;" json:"env_value" binding:"required"`
	IsSecret  bool        `gorm:"not null;" json:"is_secret" binding:"required"`
}

var dbInstance *gorm.DB

func initDB() {

	db, err := gorm.Open("sqlite3", "test.db")
	check(err)

	db.LogMode(true)

	dbInstance = db
}

func getDBInstance() *gorm.DB {
	return dbInstance
}

func migration() {
	db := getDBInstance()
	// db.DropTable(&UserInfo{}, &RoleInfo{}, &ProjectInfo{}, &ProjectEnv{})
	db.AutoMigrate(&UserInfo{}, &RoleInfo{}, &ProjectInfo{}, &ProjectEnv{})
}

func seedData(dataList ...interface{}) {

	db := getDBInstance()
	for _, v := range dataList {
		db.Create(v)
	}
}

func seed() {
	seedData(
		&UserInfo{Username: "admin", Password: encrypt("admin"), RoleID: 1},
		&RoleInfo{Name: "admin"},
		&RoleInfo{Name: "user"},
		&RoleInfo{Name: "guest"},
	)
}
