package repository

import (
	"strconv"
	"time"
)

type Video struct {
	Id            int64     `json:"id" gorm:"column:id"`
	PublisherId   int64     `json:"publisher_id" gorm:"column:publisher_id"`
	Title         string    `json:"title" gorm:"column:title"`
	VideoUrl      string    `json:"play_url" gorm:"column:video_url"`
	CoverUrl      string    `json:"cover_url" gorm:"column:cover_url"`
	FavoriteCount int64     `json:"favorite_count" gorm:"column:favorite_count"`
	CommentCount  int64     `json:"comment_count" gorm:"column:comment_count"`
	CreateDate    time.Time `json:"create_date" gorm:"column:create_date"`
}

type VideoUser struct {
	Video Video `gorm:"embedded"`
	User  User  `gorm:"embedded"`
}

type VideoDao struct{}

var videoDao = VideoDao{}

func (*VideoDao) Add(video *Video) (int64, error) {
	err := db.Table("video").Create(video).Error
	if err != nil {
		return -1, err
	}
	return video.Id, nil
}

func (*VideoDao) QueryById(id int64) (*Video, error) {
	var video Video
	err := db.Table("video").Where("id = ?", id).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (*VideoDao) QueryByPublisher(publisher int64) ([]Video, error) {
	var videos []Video
	err := db.Table("video").Where("publisher_id = ?", publisher).Order("create_date desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

//根据limit选择数据库
func (*VideoDao) Query(limit int) ([]Video, error) {
	var videos []Video
	err := db.Raw("SELECT id, publisher_id,title, video_url, cover_url, favorite_count, comment_count, create_date FROM video order by create_date desc limit " + strconv.FormatInt(int64(limit), 10)).Scan(&videos).Error
	//err := db.Table("video").Select("video.id as id, video_url, cover_url, favorite_count, comment_count, title, user.id as publisher_id, user.username as username, user.follow_count as follow_count, user.follower_count as follower_count").Joins("left join user on video.publisher_id = user.id").Scan(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (*VideoDao) QueryUserVideo(userId int64, limit int64) ([]Video, error) {
	var videos []Video
	err := db.Table("video").Where("publisher_id = ?", userId).Find(&videos).Limit(int(limit)).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (*VideoDao) Update(video *Video) error {
	return db.Table("video").Save(video).Error
}

func GetVideoDaoInstance() *VideoDao {
	return &videoDao
}
