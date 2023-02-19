package dao

import (
	"Tiktok_simple/model"
	"log"
)

// GetUserList 获取全部User对象
func GetUserList() ([]model.User, error) {
	var tableUsers []model.User
	if err := Db.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

// GetUserByUsername 根据username获得User对象
func GetUserByUsername(name string) (model.User, error) {
	tableUser := model.User{}
	if err := Db.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// GetUserById 根据user_id获得User对象
func GetUserById(id int64) (model.User, error) {
	tableUser := model.User{}
	if err := Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// InsertUser 将tableUser插入表内
func InsertUser(tableUser *model.User) bool {
	if err := Db.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
