package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"tiktok/model"
	"tiktok/service"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}
type UserResponse struct {
	Response
	User model.UserInfo `json:"user"`
}

// Register POST douyin/user/register/ 用户注册
func Register(c *gin.Context) {
	// 获取请求参数
	username := c.Query("username")
	password := c.Query("password")

	// 判断用户名是否存在
	IsExist := service.IsUserExistByUsername(username)

	if IsExist { //用户已经存在
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// 密码加盐加密
		hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "密码不符合规范"},
			})
		} else {

			//准备好userInfo,默认name为username
			userLogin := model.UserLogin{Username: username, Password: string(hashPwd)}
			userinfo := model.UserInfo{User: &userLogin, Name: username}
			println(userinfo.Name)
			if service.InsertTableUser(&userinfo) != true {
				log.Println("Insert Data Fail")
			}
			u := service.GetTableUserByUsername(username)
			token := "token" //生成token
			log.Println("注册返回的id: ", u.Id)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   u.Id,
				Token:    token,
			})
		}

	}
}

// Login POST douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	u := service.GetTableUserByUsername(username)
	// 密码对比
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Username or Password Error"},
		})
	} else {
		token := "token" //生成token
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

// UserInfo GET douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	//user_id := "1"
	// 转换为数字
	id, _ := strconv.ParseInt(user_id, 10, 64)
	println(id)
	//根据id查询user
	var u model.UserInfo
	if err := service.QueryUserInfoById(id, &u); err != nil {
		log.Println("查询失败")
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User Doesn't Exist"},
		})
	} else {
		log.Println("查询成功")

		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     u,
		})
	}
}
