package service

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"mime/multipart"

	"os"
	"strconv"
	"sync"
	"tiktok_video/repository"
	"tiktok_video/util"
	"time"
)

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func convertUser(user *repository.User, isFollow bool) *User {
	return &User{
		Id:            user.Id,
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
}

type VideoUser struct {
	repository.Video
	IsFavorite bool `json:"is_favorite"`
	User       User `json:"author"`
}

func GetVideoFeed(name string) (io.Reader, int64, time.Time, error) {
	return util.GetObjectWithSize(name)
}

func GetVideoCover(name string) (io.Reader, int64, time.Time, error) {
	return util.GetObjectWithSize(name)
}

func PublishVideo(data *multipart.FileHeader, title string, userId int64) (int64, error) {
	file, _ := data.Open()
	defer file.Close()

	fileName := strconv.FormatInt(userId, 10) + "/" + strconv.FormatInt(time.Now().Unix(), 10)

	err := util.PutVideo(fileName, file, data.Size)
	if err != nil {
		return -1, err
	}
	reader := bytes.NewBuffer(nil)
	err = ffmpeg.Input("http://localhost:8100"+util.GET_VIDEO_PATH+"?name="+fileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 5)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(reader, os.Stdout).
		Run()
	if err != nil {
		return -1, err
	}
	coverName := fileName + "-cover.jpeg"
	util.PutJpg(coverName, reader, int64(reader.Len()))

	return repository.GetVideoDaoInstance().Add(&repository.Video{
		PublisherId:   userId,
		Title:         title,
		VideoUrl:      util.OSS_SHARE_HOST + fileName,
		CoverUrl:      util.OSS_SHARE_HOST + coverName,
		FavoriteCount: 0,
		CommentCount:  0,
		CreateDate:    time.Now(),
	})
}

/*
获取视频流 30个视频
*/
func VideoList(userId int64) ([]VideoUser, error) {
	videoData, err := repository.GetVideoDaoInstance().Query(util.FEED_LIMIT) //获取dao，用来访问数据库
	if err != nil {
		return nil, err
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	var followMap, favoriteMap map[int64]struct{}
	var followErr, favoriteErr error
	earlyDate := videoData[len(videoData)-1].CreateDate
	go func() {
		//todo:这里需要relation的follow信息
		followMap, followErr = getFollowMap(userId)
		waitGroup.Done()
	}()
	go func() {
		//todo:这里需要relation的favorite信息
		favoriteMap, favoriteErr = getFavoriteMap(userId, earlyDate)
		waitGroup.Done()
	}()
	waitGroup.Wait()
	if followErr != nil {
		return nil, followErr
	}
	if favoriteErr != nil {
		return nil, favoriteErr
	}

	result := make([]VideoUser, len(videoData))
	for i, video := range videoData {
		publiser_id := video.PublisherId
		//todo:rpc调用，获取user 信息
		user := getUserInfo(publiser_id)
		raw := VideoUser{video, true, *convertUser(&user, true)}
		if _, ok := followMap[user.Id]; !ok {
			raw.User.IsFollow = false
		}
		if _, ok := favoriteMap[video.Id]; !ok {
			raw.IsFavorite = false
		}
		result[i] = raw
	}

	return result, nil
}

func getFavoriteMap(id int64, date time.Time) (map[int64]struct{}, error) {
	return nil, nil
}

func getFollowMap(id int64) (map[int64]struct{}, error) {
	return nil, nil
}

func QueryUserVideo(userId int64) ([]VideoUser, error) {
	data, err := repository.GetVideoDaoInstance().QueryUserVideo(userId, 30)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []VideoUser{}, err
	}
	favoriteMap, err := getFavoriteMap(userId, data[len(data)-1].CreateDate)
	if err != nil {
		return nil, err
	}
	result := make([]VideoUser, len(data))
	for i, video := range data {
		isFavorite := false
		if _, ok := favoriteMap[video.Id]; ok {
			isFavorite = true
		}
		//todo:根据user_id获取用户信息
		user_info := getUserInfo(userId)
		result[i] = VideoUser{
			Video:      video,
			IsFavorite: isFavorite,
			User:       *convertUser(&user_info, false),
		}
	}

	return result, nil
}
