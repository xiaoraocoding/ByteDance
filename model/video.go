package model

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title"`
}

type Video_sql struct {
	Id            int64  `json:"id,omitempty"`
	Author_id     int    `json:"author_id"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	Title         string `json:"title"`
	CreateTime    int64  `json:"create_time"`
}

type Comment struct {
	Id         int64  `json:"id"`
	Commenter  User   `json:"user"`
	Content    string `json:"content"`
	UpdateDate string `json:"create_date"`
}

type Comment_sql struct {
	Id         int64  `json:"id"`
	Comment    string `json:"comment"`
	UserId     int64  `json:"user_id"`
	UpdateTime string `json:"update_time"`
}

type Video_comment struct {
	Id        int64 `json:"id"`
	VideoId   int64 `json:"video_id"`
	CommentId int64 `json:"comment_id"`
}

func Test(c *gin.Context) {
	comment := Comment_sql{}
	comment.Comment = "这是一个测试文件"
	comment.UserId = 1

	timeStr := time.Now().Format("2006-01-02") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	comment.UpdateTime = timeStr

	Db_write.Table("comment").Create(&comment)
	c.JSON(200, gin.H{
		"msg": "ok",
	})

}
