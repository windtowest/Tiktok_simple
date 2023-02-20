package model

import (
	"errors"
	"log"
)

// UserInfo 用户信息表
type UserInfo struct {
	Id            int64      `json:"id" gorm:"id,omitempty"`
	Name          string     `json:"name" gorm:"name,omitempty"`
	FollowCount   int64      `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64      `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool       `json:"is_follow" gorm:"is_follow,omitempty"`
	User          *UserLogin `json:"-"` //用户与账号密码之间的一对一
	Videos        []*Video   `json:"-"` //用户与投稿视频的一对多
	//Follows       []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos []*Video   `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments    []*Comment `json:"-"`                                     //用户与评论的一对多
}

// UserLogin 用户登录表，和UserInfo属于一对一关系
type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

// InsertTableUser 将tableUser插入表内
func InsertTableUser(tableUser *UserInfo) bool {
	if err := DB.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// GetUserByName 通过名字查询用户
func GetUserByName(name string) (UserLogin, error) {
	var userLogin UserLogin
	if err := DB.Where("username = ?", name).First(&userLogin).Error; err != nil {
		log.Println(err.Error())
		return userLogin, err
	}
	return userLogin, nil
}

// GetTableUserByUsername 根据username获得TableUser对象
func GetTableUserByUsername(name string) (UserLogin, error) {
	tableUser := UserLogin{}
	if err := DB.Where("Username = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}
func QueryUserInfoById(userId int64, userInfo *UserInfo) error {
	DB.Where("id=?", userId).Select([]string{"id", "name", "follow_count", "follower_count", "is_follow"}).First(userInfo)
	//id为零值，说明sql执行失败
	if userInfo.Id == 0 {

		return errors.New("该用户不存在")
	}
	return nil
}
