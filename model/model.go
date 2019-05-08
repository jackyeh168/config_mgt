package model

import (
	"auth/util"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserInfo struct {
	gorm.Model
	Username string        `gorm:"type:varchar(16);not null;unique_index" json:"username" binding:"required"`
	Password string        `gorm:"type:varchar(128);not null;" json:"password"`
	RoleID   uint          `sql:"type:int unsigned" gorm:"not null;" json:"role_id"`
	Role     RoleInfo      `gorm:"foreignkey:RoleID" json:"-"`
	Projects []ProjectInfo `gorm:"many2many:user_projects;" json:"-"`
}

type RoleInfo struct {
	gorm.Model
	Name       string `gorm:"type:varchar(16);not null;unique_index" json:"role_name"`
	Permission string `gorm:"type:varchar(255);" json:"permission"`
}

type ProjectInfo struct {
	gorm.Model
	Name  string     `gorm:"type:varchar(32);not null;unique_index" json:"project_name" binding:"required"`
	Users []UserInfo `gorm:"many2many:user_projects;" json:"-"`
}

type UserProject struct {
	UserInfoID    uint `gorm:"not null;unique_index:idx_user_proj" json:"user_info_id" binding:"required"`
	ProjectInfoID uint `gorm:"not null;unique_index:idx_user_proj" json:"project_info_id" binding:"required"`
}

type ProjectEnv struct {
	gorm.Model
	ProjectID uint        `sql:"type:int unsigned" gorm:"not null;unique_index:idx_proj_env" json:"project_id" binding:"required"`
	Project   ProjectInfo `gorm:"foreignkey:ProjectID" json:"-"`
	EnvKey    string      `gorm:"type:varchar(32);not null;unique_index:idx_proj_env" json:"env_key" binding:"required"`
	EnvValue  string      `gorm:"type:varchar(128);not null;" json:"env_value" binding:"required"`
}

var dbInstance *gorm.DB

func InitDB() {

	db, err := gorm.Open("sqlite3", "test.db")
	util.Check(err)

	db.LogMode(true)

	dbInstance = db
}

func GetDBInstance() *gorm.DB {
	return dbInstance
}

func Migration() {
	db := GetDBInstance()
	db.DropTable(&UserInfo{}, &RoleInfo{}, &ProjectInfo{}, &ProjectEnv{}, &UserProject{})
	db.AutoMigrate(&UserInfo{}, &RoleInfo{}, &ProjectInfo{}, &ProjectEnv{}, &UserProject{})
}

func seedData(dataList ...interface{}) {
	db := GetDBInstance()
	for _, v := range dataList {
		db.Create(v)
	}
}

func Seed() {
	seedData(
		&UserInfo{Username: "admin", Password: util.Encrypt("admin"), RoleID: 1},
		&RoleInfo{Name: "admin", Permission: `{"user":["create","read","update","delete"],"project":["create","read","update","delete"],"env":["create","read","update","delete"]}`},
		&RoleInfo{Name: "leader", Permission: `{"user":[],"project":["read"],"env":["create","read","update","delete"]}`},
		&RoleInfo{Name: "guest", Permission: `{"user":[],"project":["read"],"env":["read"]}`},
	)
}
