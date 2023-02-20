package main

import (
	"github.com/gin-gonic/gin"
	"tiktok/model"
)

func main() {
	// 初始化数据库服务
	db := model.InitDB()
	defer db.Close()
	// 创建一个服务
	r := gin.Default()
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
