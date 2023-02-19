package service

import (
	"Tiktok_simple/kitex_gen/relation_gorm"
	"Tiktok_simple/model"
)

type UserService interface {
	/*
		个人使用
	*/
	// GetUserList 获得全部User对象
	GetUserList() []model.User

	// GetUserByUsername 根据username获得User对象
	GetUserByUsername(name string) model.User

	// InsertUser 将tableUser插入表内
	InsertUser(tableUser *model.User) bool
	/*
		他人使用
	*/
	// GetUserById 未登录情况下,根据user_id获得User对象
	GetUserById(id int64) (relation_gorm.User, error)

	// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
	GetUserByIdWithCurId(id int64, curId int64) (*relation_gorm.User, error)

	// 根据token返回id
	// 接口:auth中间件,解析完token,将userid放入context
	//(调用方法:直接在context内拿参数"userId"的值)	fmt.Printf("userInfo: %v\n", c.GetString("userId"))
}

// User 最终封装后,controller返回的User结构体
type User struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
}
