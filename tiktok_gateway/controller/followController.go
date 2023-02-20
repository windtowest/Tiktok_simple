package controller

import (
	"TikTok/config"
	"Tiktok_simple/kitex_gen/relation_gorm"
	"context"
	"github.com/cloudwego/kitex/client"

	"Tiktok_simple/kitex_gen/relation_gorm/userservice"
	"fmt"
	"github.com/gin-gonic/gin"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// RelationActionResp 关注和取消关注需要返回结构。
type RelationActionResp struct {
	Response
}

// FollowingResp 获取关注列表需要返回的结构。
type FollowingResp struct {
	Response
	UserList []relation_gorm.User `json:"user_list,omitempty"`
}

// FollowersResp 获取粉丝列表需要返回的结构。
type FollowersResp struct {
	Response
	// 必须大写，才能序列化
	UserList []relation_gorm.User `json:"user_list,omitempty"`
}

var (
	followClient      userservice.Client //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once          //限定该service对象为单例，节约内存。
)

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() userservice.Client {
	followServiceOnce.Do(
		func() {
			r, err := etcd.NewEtcdResolver([]string{config.EtcdAddr})
			if err != nil {
				log.Fatal(err)
			}
			followClient, err = userservice.NewClient("relation", client.WithResolver(r))
			if err != nil {
				log.Fatal(err)
			}
		})
	return followClient
}

// RelationAction 处理关注和取消关注请求。
func RelationAction(c *gin.Context) {

	userId, err1 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)

	if nil != err1 || nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误",
			},
		})
		return
	}
	req := &relation_gorm.ActionRequest{
		UserId:     userId,
		ToUserId:   toUserId,
		ActionType: int32(actionType),
	}
	client := NewFSIInstance()
	resp, err := client.Action(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, resp)
	return
}

// GetFollowing 处理获取关注列表请求。
func GetFollowing(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// 用户id解析出错。
	if nil != err {
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误。",
			},
			UserList: nil,
		})
		return
	}
	req := &relation_gorm.FollowListRequest{
		UserId: userId,
	}
	// 正常获取关注列表
	client := NewFSIInstance()
	resp, err := client.FollowList(context.Background(), req)

	// 获取关注列表时出错。
	if err != nil {
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表时出错。",
			},
			UserList: nil,
		})
		return
	}
	// 成功获取到关注列表。
	log.Println("获取关注列表成功。")
	c.JSON(http.StatusOK, resp)
}

// GetFollowers 处理获取关注列表请求
func GetFollowers(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// 用户id解析出错。
	if nil != err {
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误。",
			},
			UserList: nil,
		})
		return
	}
	req := &relation_gorm.FollowerListRequest{
		UserId: userId,
	}
	// 正常获取粉丝列表
	client := NewFSIInstance()
	resp, err := client.FollowList(context.Background(), req)
	// 获取关注列表时出错。
	if err != nil {
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表时出错。",
			},
			UserList: nil,
		})
		return
	}
	// 成功获取到粉丝列表。
	//log.Println("获取粉丝列表成功。")
	c.JSON(http.StatusOK, resp)
}

// GetFollowers 处理获取关注列表请求
func GetFriendList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// 用户id解析出错。
	if nil != err {
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误。",
			},
			UserList: nil,
		})
		return
	}
	req := &relation_gorm.FollowerListRequest{
		UserId: userId,
	}
	// 正常获取好友列表
	client := NewFSIInstance()
	resp, err := client.FriendList(context.Background(), req)
	// 获取好友列表时出错。
	if err != nil {
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取好友列表时出错。",
			},
			UserList: nil,
		})
		return
	}
	// 成功获取到粉丝列表。
	//log.Println("获取粉丝列表成功。")
	c.JSON(http.StatusOK, resp)
}
