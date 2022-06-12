package service

import (
	"ByteDance/config"
	"ByteDance/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {

	uid,_ := c.Get("user_id")
	userid := config.GetInterfaceToString(uid)
	video_id := c.Query("video_id")
	action_type := c.Query("action_type") //1是点赞，2是取消点赞

	res,_:= strconv.Atoi(action_type)

	if res  == 1 {
		video := model.Video_sql{}
		v_id,_ := strconv.Atoi(video_id)
		video.Id = int64(v_id)
		res := model.Rdb.SMembers(model.Ctx,video_id)
		res_list,_ := res.Result()
		if len(res_list)!=0{
			model.Rdb.Incr(model.Ctx,video_id) //视频点赞数加1
			sum := model.Rdb.Get(model.Ctx,video_id)
			s_res,_ := sum.Result()
			s,_ := strconv.Atoi(s_res)
			model.Db_write.Table("videos").Model(&video).Update("favorite_count",s)
		}else{
			//不是当天发布的视频的情况
			model.Db_write.Table("videos").Where("video_id = ?",video.Id).UpdateColumn("favorite_count",gorm.Expr("favorite_count + ?",1))
		}
		model.Rdb.SAdd(model.Ctx,userid,video_id) //点赞集合加1
		like := model.Like{}
		//得先判断之前有没有点赞过
		model.Db_write.Table("like").Where("user_id = ? AND video_id = ?",userid,video_id).First(&like)
		if like.Id==0{
			id,_ := strconv.Atoi(userid)
			like.UserId = int64(id)
			like.VideoId = int64(v_id)
			like.ActionType = 1
			model.Db_write.Table("like").Create(&like)
		}else{
			model.Db_write.Table("like").Where("user_id = ? AND video_id = ?",userid,video_id).Update("action_type",1)
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg: "",
		})
	}else {
		video := model.Video_sql{}
		v_id,_ := strconv.Atoi(video_id)
		video.Id = int64(v_id)
		res := model.Rdb.SMembers(model.Ctx,video_id)
		res_list,_ := res.Result()
		if len(res_list)!=0{
			model.Rdb.IncrBy(model.Ctx,video_id,-1)
			sum := model.Rdb.Get(model.Ctx,video_id)
			s_res,_ := sum.Result()
			s,_ := strconv.Atoi(s_res)
			model.Db_write.Table("videos").Model(&video).Update("favorite_count",s)
		}else{
			model.Db_write.Table("videos").Where("video_id = ?",video.Id).UpdateColumn("favorite_count",gorm.Expr("favorite_count - ?",1))
		}
		model.Rdb.SRem(model.Ctx,userid,video_id)
		like := model.Like{}
		id,_ := strconv.Atoi(userid)
		like.Id = int64(id)
		like.VideoId=int64(v_id)
		model.Db_write.Table("like").Where("user_id = ? AND video_id = ?",like.UserId,like.VideoId).Update("action_type",2)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg: "",
		})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	uid,_ := c.Get("user_id")
	user_id := config.GetInterfaceToString(uid)
	videoSqlList := []model.Video_sql{}

	res := model.Rdb.SMembers(model.Ctx,user_id)
	res_list,_ := res.Result()
	fmt.Println(res_list)   //这就是用户关注的视频列表，里面的数据就是视频id
	if len(res_list)!=0{
		model.Db_write.Table("videos").Where("id in (?)", res_list).Find(&videoSqlList)
	}

	videoList := make([]model.Video, len(videoSqlList))
	for i, _ := range videoSqlList {
		videoList[i].Id = videoSqlList[i].Id
		videoList[i].Title = videoSqlList[i].Title
		videoList[i].CommentCount = videoSqlList[i].CommentCount
		videoList[i].FavoriteCount = videoSqlList[i].FavoriteCount
		videoList[i].PlayUrl = videoSqlList[i].PlayUrl
		videoList[i].CoverUrl = videoSqlList[i].CoverUrl
		videoList[i].IsFavorite = true
		user := model.User{}
		model.Db_write.Table("user").Where("id = ?", videoSqlList[i].Author_id).First(&user)
		videoList[i].Author = user
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
