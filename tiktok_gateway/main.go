package main

import (
	"github.com/gin-gonic/gin"
)

// 如果启动有问题，大概是你的IP地址出现变化，需要在项目依赖的服务器中配置安全组
func main() {
	//gin
	r := gin.Default()
	initRouter(r)
	//pprof
	//pprof.Register(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
