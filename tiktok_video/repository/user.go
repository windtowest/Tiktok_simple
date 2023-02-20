package repository

type User struct {
	Id            int64  `gorm:"column:id"`
	Username      string `gorm:"column:username"`
	Password      string `gorm:"column:password"`
	FollowCount   int64  `gorm:"column:follow_count"`
	FollowerCount int64  `gorm:"column:follower_count"`
}
