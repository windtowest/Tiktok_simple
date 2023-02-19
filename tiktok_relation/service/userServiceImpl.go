package service

import (
	"Tiktok_simple/config"
	"Tiktok_simple/dao"
	"Tiktok_simple/kitex_gen/relation_gorm"
	"Tiktok_simple/model"
	"crypto/hmac"
	"crypto/sha256"

	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

type UserServiceImpl struct {
	FollowService
}

// GetUserList 获得全部User对象
func (usi *UserServiceImpl) GetUserList() []model.User {
	tableUsers, err := dao.GetUserList()
	if err != nil {
		log.Println("Err:", err.Error())
		return tableUsers
	}
	return tableUsers
}

// GetUserByUsername 根据username获得User对象
func (usi *UserServiceImpl) GetUserByUsername(name string) model.User {
	tableUser, err := dao.GetUserByUsername(name)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	log.Println("Query User Success")
	return tableUser
}

// InsertUser 将tableUser插入表内
func (usi *UserServiceImpl) InsertUser(tableUser *model.User) bool {
	flag := dao.InsertUser(tableUser)
	if flag == false {
		log.Println("插入失败")
		return false
	}
	return true
}

// GetUserById 未登录情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserById(id int64) (relation_gorm.User, error) {
	user := relation_gorm.User{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	tableUser, err := dao.GetUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")
	followCount, _ := usi.GetFollowingCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	followerCount, _ := usi.GetFollowerCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}

	user = relation_gorm.User{
		Id:            id,
		Name:          tableUser.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      false,
	}
	return user, nil
}

// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserByIdWithCurId(id int64, curId int64) (*relation_gorm.User, error) {
	user := relation_gorm.User{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	tableUser, err := dao.GetUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return &user, err
	}
	log.Println("Query User Success")
	followCount, err := usi.GetFollowingCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	followerCount, err := usi.GetFollowerCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	isfollow, err := usi.IsFollowing(curId, id)
	if err != nil {
		log.Println("Err:", err.Error())
	}

	user = relation_gorm.User{
		Id:            id,
		Name:          tableUser.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isfollow,
	}
	return &user, nil
}

// GenerateToken 根据username生成一个token
func GenerateToken(username string) string {
	u := UserService.GetUserByUsername(new(UserServiceImpl), username)
	fmt.Printf("generatetoken: %v\n", u)
	token := NewToken(u)
	println(token)
	return token
}

// NewToken 根据信息创建token
func NewToken(u model.User) string {
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id64 := u.Id
	fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
	claims := jwt.StandardClaims{
		Audience:  u.Name,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(id64, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	var jwtSecret = []byte(config.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		token = "Bearer " + token
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

// EnCoder 密码加密
func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Result: " + sha)
	return sha
}
