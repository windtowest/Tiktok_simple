package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok_video/service"
)

type PublishListResponse struct {
	BaseResponse
	VideoList []service.VideoUser `json:"video_list"`
	NextTime  int64               `json:"next_time"`
}

func Publish(c *gin.Context) {

	tokenString := c.Query("token")
	//token, _ := service.GetAuthInstance().ParseTokenString(tokenString)
	//todo:解析jwt 获取id
	userId := parser(tokenString)
	//for key, value := range token.Claims.(jwt.MapClaims) {
	//	if key == "identity" {
	//		userId = int64(int(value.(float64)))
	//	}
	//}

	//c.JSON(http.StatusOK, service.GetUserId(c))
	file, _ := c.FormFile("data")
	title := c.PostForm("title")
	_, err := service.PublishVideo(file, title, userId)
	if err != nil {
		buildError(c, err.Error())
		return
	}
	buildSuccess(c)
}

func PublishList(c *gin.Context) {
	tokenString := c.Query("token")
	//todo:parser token string
	token_id := parser(tokenString)
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//todo:这个鉴权也可以移到网关去做
	if token_id != userId {
		buildError(c, err.Error())
		return
	}
	if err != nil {
		buildError(c, err.Error())
		return
	}
	data, err := service.QueryUserVideo(userId)
	fmt.Println("publishlist")
	fmt.Println(data)
	if err != nil {
		buildError(c, err.Error())
	}
	if len(data) == 0 {
		c.JSON(http.StatusOK, PublishListResponse{
			BaseResponse: BaseResponse{
				StatusCode: 0,
				StatusMsg:  "success",
			},
		})
	}
	c.JSON(http.StatusOK, PublishListResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: data,
		NextTime:  data[len(data)-1].CreateDate.UnixNano(),
	})
}

func parser(tokenString string) int64 {
	return 0
}
