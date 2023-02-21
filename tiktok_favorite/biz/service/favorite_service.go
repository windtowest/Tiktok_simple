package service

import (
	"Tiktok_simple/tiktok_favorite/biz/dao/db"
	// "Tiktok_simple/tiktok_favorite/biz/model/basic/feed"
	favorite "Tiktok_simple/tiktok_favorite/biz/model/interact/favorite"
	"Tiktok_simple/tiktok_favorite/pkg/constants"
	"Tiktok_simple/tiktok_favorite/pkg/errno"
	"context"

	// feed_service "Tiktok_simple/biz/service/feed"

	"github.com/cloudwego/hertz/pkg/app"
)

type FavoriteService struct {
	ctx context.Context
	c   *app.RequestContext
}

// new FavoriteService
func NewFavoriteService(ctx context.Context, c *app.RequestContext) *FavoriteService {
	return &FavoriteService{ctx: ctx, c: c}
}

// like action, include like and unlike
// request parameters:
// string token = 1;       // 用户鉴权token
// int64 to_user_id = 2;   // 对方用户id
// int32 action_type = 3;  // 1-点赞，2-取消点赞
func (r *FavoriteService) FavoriteAction(req *favorite.DouyinFavoriteActionRequest) (flag bool, err error) {
	// 颁发和验证token的行为均交给jwt处理，当发送到handler层时，默认已通过验证
	// 只需要检查参数VideoID的合法性

	_, err = db.CheckVideoExistById(req.VideoId)
	if err != nil {
		return false, err
	}
	if req.ActionType != constants.FavoriteActionType && req.ActionType != constants.UnFavoriteActionType {
		return false, errno.ParamErr
	}
	// 获取current_user_id
	current_user_id, _ := r.c.Get("")
	// // 不准自己关注自己
	// if req.ToUserId == current_user_id.(int64) {
	// 	return false, errno.ParamErr
	// }
	new_favorite_relation := &db.Favorites{
		UserId:  current_user_id.(int64),
		VideoId: req.VideoId,
	}
	// 请求参数校验完毕，检查favorite表中是否已经存在这两者的关系
	favorite_exist, _ := db.CheckFavoriteRelationExist(new_favorite_relation)
	if req.ActionType == constants.FavoriteActionType {
		if favorite_exist {
			return false, errno.FavoriteRelationAlreadyExistErr
		}
		flag, err = db.AddNewFavorite(new_favorite_relation)
	} else {
		if !favorite_exist {
			return false, errno.FavoriteRelationNotExistErr
		}
		flag, err = db.DeleteFavorite(new_favorite_relation)
	}
	return flag, err
}

// 获取用户点赞的所有视频列表，需要注意的是这里的token是客户端当前用户，而user_id可以是任意用户//zheli
// request parameters:
// string token;       // 用户鉴权token
// int64  user_id;     // 用户id
func (r *FavoriteService) GetFavoriteList(req *favorite.DouyinFavoriteListRequest) (favoritelist []*favorite.Video, err error) {
	query_user_id := req.UserId
	_, err = db.CheckUserExistById(query_user_id)

	if err != nil {
		return nil, err
	}
	// 获取current_user_id
	current_user_id, _ := r.c.Get("")

	video_id_list, err := db.GetFavoriteIdList(query_user_id)

	dbVideos, err := db.GetVideoListByVideoIDList(video_id_list)
	var videos []*feed.Video
	f := feed_service.NewFeedService(r.ctx, r.c)
	err = f.CopyVideos(&videos, &dbVideos, current_user_id.(int64))
	for _, item := range videos {
		video := &favorite.Video{
			Id: item.Id,
			Author: favorite.User{
				Id:            item.Author.Id,
				Name:          item.Author.Name,
				FollowCount:   item.Author.FollowCount,
				FollowerCount: item.Author.FollowerCount,
				IsFollow:      item.Acthor.IsFollow,
			},
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: item.FavoriteCount,
			CommentCount:  item.CommentCount,
			IsFavorite:    item.IsFavorite,
			Title:         item.Title,
		}
		favoritelist = append(favoritelist, video)
	}
	return favoritelist, err
}
