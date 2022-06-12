package service

import (
	"ByteDance/config"
	"ByteDance/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userid := config.GetInterfaceToString(uid)
	last_time := c.Query("last_time")
	if last_time == "" {
		nowTime := time.Now().Unix()
		last_time = strconv.Itoa(int(nowTime))
	}
	fmt.Println("last_time:", last_time)
	video_sql := []model.Video_sql{}

	model.Db_write.Table("videos").Limit(5).Order("create_time desc").Where("create_time < ?", last_time).Find(&video_sql)
	fmt.Println("video_sql数", len(video_sql))
	next_time := video_sql[len(video_sql)-1].CreateTime

	video := make([]model.Video, len(video_sql))
	for i := 0; i < len(video_sql); i++ {
		video[i].Id = video_sql[i].Id
		author := model.User{}
		model.Db_write.Table("user").Where("id = ?", video_sql[i].Author_id).First(&author)
		user_follow_list := userid + "/follow/list"
		res := model.Rdb.SIsMember(model.Ctx, user_follow_list, author.Id)
		isFo, _ := res.Result()
		if isFo {
			author.IsFollow = true
		} else {
			author.IsFollow = false
		}
		video[i].Title = video_sql[i].Title
		video[i].CommentCount = video_sql[i].CommentCount
		video[i].FavoriteCount = video_sql[i].FavoriteCount

		video[i].PlayUrl = video_sql[i].PlayUrl
		video[i].CoverUrl = video_sql[i].CoverUrl
		video[i].Author = author
		// 从redis 中判断
		res2 := model.Rdb.SIsMember(model.Ctx, userid, strconv.Itoa(int(video[i].Id)))
		isFa, _ := res2.Result()
		if isFa {
			video[i].IsFavorite = true
		} else {
			video[i].IsFavorite = false
		}
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: video,
		NextTime:  next_time,
	})
}
