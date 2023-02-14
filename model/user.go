package model

// User 对应数据库User表结构的结构体
type User struct {
	Id       int64
	Name     string
	Password string
}

// TableName 修改表名映射
func (tableUser User) TableName() string {
	return "users"
}
