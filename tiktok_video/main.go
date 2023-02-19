package tiktok_video

import (
	"github.com/gin-gonic/gin"
	"os"
	"tiktok_video/repository"
	"tiktok_video/util"
)

func main() {
	err := Init()
	if err != nil {
		os.Exit(1)
	}
	r := gin.Default()
	initRouter(r)
	err = r.Run(":8100")
	if err != nil {
		return
	}
}

func Init() error {
	if err := repository.Init(); err != nil {
		return err
	}
	if err := util.InitOSS(); err != nil {
		return err
	}
	return nil
}
