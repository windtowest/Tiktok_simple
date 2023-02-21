package db

import (
	"Tiktok_simple/tiktok_favorite/pkg/constants"
	"Tiktok_simple/tiktok_favorite/pkg/errno"
	"time"

	"gorm.io/gorm"
)

type Favorites struct {
	ID        int64          `json:"id"`
	UserId    int64          `json:"user_id"`
	VideoId   int64          `json:"video_id"`
	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (f *Favorites) TableName() string {
	return constants.FavoritesTableName
}

func AddNewFavorite(favorite *Favorites) (bool, error) {
	err := DB.Create(favorite).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteFavorite(favorite *Favorites) (bool, error) {
	err := DB.Where("video_id = ? AND user_id = ?", favorite.VideoId, favorite.UserId).Delete(favorite).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func QueryFavoriteExist(video_id int64, user_id int64) (bool, error) {
	var sum int64
	err := DB.Model(&Favorites{}).Where("video_id = ? AND user_id = ?", video_id, user_id).Count(&sum).Error

	if err != nil {
		return false, err
	}
	if sum == 0 {
		return false, nil
	}
	return true, nil
}

func QueryTotalFavoritedByAuthorID(author_id int64) (int64, error) {
	var sum int64
	err := DB.Table("likes").Joins("JOIN videos ON likes.video_id = videos.id").
		Where("videos.author_id = ?", author_id).Count(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// 查询视频的点赞数量
func GetFavoriteCount(video_id int64) (int64, error) {
	var count int64
	err := DB.Model(&Favorites{}).Where("video_id = ?", video_id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 获得 user_id 点赞的视频的 video_id
func GetFavoriteIdList(user_id int64) ([]int64, error) {
	var favorite_actions []Favorites
	err := DB.Where("user_id = ?", user_id).Find(&favorite_actions).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range favorite_actions {
		result = append(result, v.VideoId)
	}
	return result, nil
}

// 获得 user点赞的视频视频数量
func GetFavoriteCountByUserID(user_id int64) (int64, error) {
	var count int64
	err := DB.Model(&Favorites{}).Where("user_id = ?", user_id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 获得 video_id 所有点赞的人的 id
func GetFavoriterIdList(video_id int64) ([]int64, error) {
	var favorite_actions []Favorites
	err := DB.Where("video_id = ?", video_id).Find(&favorite_actions).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range favorite_actions {
		result = append(result, v.UserId)
	}
	return result, nil
}

func CheckFavoriteRelationExist(favorite *Favorites) (bool, error) {
	err := DB.Where("video_id = ? AND user_id = ?", favorite.VideoId, favorite.UserId).Find(&favorite).Error
	if err != nil {
		return false, err
	}
	// find未找到符合条件的数据会返回空结构体，ID = 0
	if favorite.ID == 0 {
		err := errno.FavoriteRelationNotExistErr
		return false, err
	}
	return true, nil
}
