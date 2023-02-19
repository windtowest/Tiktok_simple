package tiktok_video

import (
	"github.com/gin-gonic/gin"
	"tiktok_video/controller"
)

func initRouter(r *gin.Engine) {
	douyinGroup := r.Group("/douyin")
	{
		//视频流获取
		douyinGroup.GET("/feed/", controller.GetFeedList)
		douyinGroup.GET("/feed/get", controller.GetFeed)
		douyinGroup.GET("/feed/get/cover", controller.GetCover)

		publishGroup := douyinGroup.Group("/publish")
		{
			publishGroup.POST("/action/", controller.Publish)
			//视频发布列表
			//todo:需要先鉴权
			publishGroup.GET("/list/", controller.PublishList)
		}
	}
}
