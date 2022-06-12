package service

import (
	"ByteDance/config"
	"ByteDance/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	comment_text := c.Query("comment_text")
	comment_id := c.Query("comment_id")
	uid,_ := c.Get("user_id")

	user_id := config.GetInterfaceToString(uid)
	user_i,_ := strconv.Atoi(user_id)

	action_t,_ := strconv.Atoi(action_type)
	comm_id,_ := strconv.Atoi(comment_id)
	video_i,_ := strconv.Atoi(video_id)

	if action_t == 1 { //此时就是发布评论
		comment := model.Comment_sql{}
		comment.Comment = comment_text
		comment.UserId = int64(user_i)

		timeStr:=time.Now().Format("2006-01-02")  //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
		comment.UpdateTime = timeStr
		model.Db_write.Table("comment").Create(&comment)

		video := model.Video_comment{}
		video.VideoId = int64(video_i)
		video.CommentId = int64(comm_id)

		user := model.User{}
		model.Db_write.Table("user").Where("id = ?",user_i).First(&user)

		newComment := model.Comment_sql{}
		model.Db_write.Table("comment").Where("comment = ? and user_id = ?",comment_text,user_id).First(&newComment)

		video_comment := model.Video_comment{
			VideoId: int64(video_i),
			CommentId: newComment.Id,
		}
		model.Db_write.Table("video_comment").Create(&video_comment)
		commentModel := model.Comment{
			Id: newComment.Id,
			Commenter: user,
			Content: comment_text,
			UpdateDate: timeStr,
		}
		comment_list := []model.Comment{commentModel}


		c.JSON(http.StatusOK,CommentListResponse{
			Response:Response{
				StatusCode: 0,
				StatusMsg: "",
			},
			CommentList: comment_list,
		})
	}else if action_t == 2 {
		comment := model.Comment_sql{}

		model.Db_write.Table("comment").Where("id = ?",comment_id).Delete(&comment)

		comment_model := model.Video_comment{}

		model.Db_write.Table("video_comment").Where("comment_id = ?",comment_id).Delete(&comment_model)

		c.JSON(http.StatusOK,CommentListResponse{
			Response:Response{
				StatusCode: 0,
				StatusMsg: "",
			},
		})


	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	video_id := c.Query("video_id")


	var commentId_all []int
	model.Db_write.Table("video_comment").Select("comment_id").Where("video_id = ?",video_id).Scan(&commentId_all)

     var commentlist []model.Comment
	for i := 0 ; i < len(commentId_all) ; i ++ {
		comment := model.Comment{}
		comment.Id = int64(commentId_all[i])
		var user_id int
		model.Db_write.Table("comment").Select("user_id").Where("id = ?",commentId_all[i]).Scan(&user_id)
		var user model.User
		model.Db_write.Table("user").Where("id = ?",user_id).First(&user)
       var time string
		model.Db_write.Table("comment").Select("update_time").Where("id = ?",commentId_all[i]).Scan(&time)
		var text string
		model.Db_write.Table("comment").Select("comment").Where("id = ?",commentId_all[i]).Scan(&text)

        comment.Commenter = user
		comment.Content = text
		comment.UpdateDate = time
		commentlist = append(commentlist,comment)
	}



	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentlist,

	})
}