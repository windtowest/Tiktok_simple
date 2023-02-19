package tiktok_relation

import (
	relation_gorm "Tiktok_simple/kitex_gen/relation_gorm"
	"Tiktok_simple/service"
	"Tiktok_simple/util"
	"context"
	"log"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Action implements the UserServiceImpl interface.
func (s *UserServiceImpl) Action(ctx context.Context, req *relation_gorm.ActionRequest) (resp *relation_gorm.ActionResponse, err error) {

	userId := util.GetUserIdByToken(req.Token)
	// 正常处理
	fsi := service.NewFSIInstance()
	switch {
	// 关注
	case 1 == req.ActionType:
		go fsi.AddFollowRelation(userId, req.ToUserId)
	// 取关
	case 2 == req.ActionType:
		go fsi.DeleteFollowRelation(userId, req.ToUserId)
	}
	log.Println("关注、取关成功。")
	resp = new(relation_gorm.ActionResponse)

	return
}

// FollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) FollowList(ctx context.Context, req *relation_gorm.FollowListRequest) (resp *relation_gorm.FollowListResponse, err error) {
	// 正常获取关注列表
	fsi := service.NewFSIInstance()
	users, err := fsi.GetFollowing(req.UserId)

	resp = new(relation_gorm.FollowListResponse)
	// 获取关注列表时出错。
	if err != nil {
		resp.UserList = nil
		resp.StatusCode = relation_gorm.Code_DBErr
		resp.StatusMsg = "获取关注列表时出错。"
		return
	}
	// 成功获取到关注列表。
	log.Println("获取关注列表成功。")
	resp.UserList = users
	return
}

// FollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) FollowerList(ctx context.Context, req *relation_gorm.FollowerListRequest) (resp *relation_gorm.FollowerListResponse, err error) {
	// 正常获取粉丝列表
	fsi := service.NewFSIInstance()
	users, err := fsi.GetFollowers(req.UserId)
	resp = new(relation_gorm.FollowerListResponse)
	// 获取粉丝列表时出错。
	if err != nil {
		resp.UserList = nil
		resp.StatusCode = relation_gorm.Code_DBErr
		resp.StatusMsg = "获取关注列表时出错。"
		return
	}
	// 成功获取到粉丝列表。
	log.Println("获取关注列表成功。")
	resp.UserList = users
	return
}

// FriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) FriendList(ctx context.Context, req *relation_gorm.FriendListRequest) (resp *relation_gorm.FriendListResponse, err error) {
	fsi := service.NewFSIInstance()
	followings, err := fsi.GetFollowing(req.UserId)
	followers, err := fsi.GetFollowers(req.UserId)

	resp = new(relation_gorm.FriendListResponse)
	if err != nil {
		resp.UserList = nil
		resp.StatusCode = relation_gorm.Code_DBErr
		resp.StatusMsg = "获取好友列表时出错。"
		return
	}

	friendUsers := make([]*relation_gorm.FriendUser, 0)
	tmp := make(map[int64]interface{})
	for _, follower := range followers {
		tmp[follower.Id] = nil
	}
	for _, following := range followings {
		if _, ok := tmp[following.Id]; ok {
			var friend relation_gorm.FriendUser
			friend.Id = following.Id
			friend.IsFollow = following.IsFollow
			friend.FollowerCount = following.FollowerCount
			friend.FollowCount = following.FollowCount
			friendUsers = append(friendUsers, &friend)
		}
	}
	resp.UserList = friendUsers
	return
}
