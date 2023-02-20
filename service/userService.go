package service

import (
	"log"
	"tiktok/model"
)

// IsUserExistByUsername 查询是否已存在对应的账号
func IsUserExistByUsername(username string) bool {

	tableUser, err := model.GetUserByName(username)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return false
	}
	log.Println("Query User Success")
	if tableUser.Id == 0 { // 不存在
		return false
	}
	return true
}

// InsertTableUser 插入新用户
func InsertTableUser(userinfo *model.UserInfo) bool {
	success := model.InsertTableUser(userinfo)
	if success == false {
		log.Println("插入新用户失败")
		return false
	}
	return true
}

// GetTableUserByUsername 根据username获得TableUser对象
func GetTableUserByUsername(name string) model.UserLogin {
	tableUser, err := model.GetTableUserByUsername(name)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	log.Println("Query User Success")
	return tableUser
}
func QueryUserInfoById(userId int64, userInfo *model.UserInfo) error {
	err := model.QueryUserInfoById(userId, userInfo)
	if err != nil {
		return err
	}
	return nil
}
