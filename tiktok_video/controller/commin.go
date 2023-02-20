package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func buildError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, BaseResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func buildSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}
