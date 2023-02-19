package service

import "tiktok_video/repository"

func getUserInfo(user_id int64) repository.User {
	return repository.User{0, "0", "", 0, 0}
}
