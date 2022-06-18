package controller

import (
	"net/http"
	"strconv"

	"gitee.com/Whitroom/imitate-tiktok/common"
	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

// 如果出现token 则不会出现自己的视频
func Feed(ctx *gin.Context) {
	db := sql.GetSession()

	var latestTime, nextTime int64
	token := ctx.Query("token")
	latestTime_ := ctx.Query("latest_time")
	if latestTime_ != "" {
		latestTime, _ = strconv.ParseInt(latestTime_, 10, 64)
	} else {
		latestTime = 0
	}
	var userID uint
	if token != "" {
		var err error
		userID, err = middlewares.Parse(ctx, token)
		if err != nil {
			return
		}

	} else {
		userID = 0
	}
	videos := crud.GetVideos(db, latestTime, userID)
	responseVideos := common.VideosModelChange(db, userID, videos)
	if len(videos)-1 < 0 {
		nextTime = 0
	} else {
		nextTime = videos[len(videos)-1].CreatedAt.Unix()
	}
	ctx.JSON(http.StatusOK, response.FeedResponse{
		Response: response.Response{
			StatusCode: response.SUCCESS,
			StatusMsg:  "获取成功",
		},
		VideoList: responseVideos,
		NextTime:  nextTime,
	})
}
