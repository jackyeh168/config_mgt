package model

import "auth/util"

func GetUsers(users []UserInfo) error {
	db := GetDBInstance()
	return db.Find(&users).Error
}

func SaveUser(user UserInfo) bool {
	if len(user.Password) < 2 || user.Password[:2] != "$2" {
		user.Password = util.Encrypt(user.Password)
	}

	db := GetDBInstance()

	if db.Where(&user).Find(&user).RecordNotFound() {
		util.Check(db.Create(&user).Error)
		return true
	}
	return false
}

func SaveUsers(users []UserInfo) bool {
	res := true
	for _, user := range users {
		res = SaveUser(user)
	}
	return res
}

func UpdateUser(user UserInfo) error {
	db := GetDBInstance()
	var olduser UserInfo
	olduser.ID = user.ID
	db.First(&olduser)

	if user.Password != "" {
		olduser.Password = util.Encrypt(user.Password)
	}

	olduser.Username = user.Username
	olduser.RoleID = user.RoleID
	return db.Save(&olduser).Error
}

func DeleteUser(user UserInfo) error {
	db := GetDBInstance()
	return db.Unscoped().Delete(&user).Error
}
