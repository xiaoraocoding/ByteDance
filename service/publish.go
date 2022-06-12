package service

import (
	"ByteDance/config"
	"ByteDance/model"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	file, err := c.FormFile("data")
	if err != nil {
		fmt.Println("form failed ....", err)
	}
	fmt.Println("file.filename:",file.Filename)
	user_id,_ := c.Get("user_id")
	uid := config.GetInterfaceToString(user_id)
	author_id,_:=strconv.Atoi(uid)
	finalFilename:=fmt.Sprintf("%d_%s",author_id,file.Filename)
	local_filename := "./public/" +  finalFilename
	filenameall:=path.Base(file.Filename)
	filesuffix:=path.Ext(file.Filename)
	covername:=filenameall[0:len(filenameall)-len(filesuffix)]+".jpg"
	finalCovername:=fmt.Sprintf("%d_%s",author_id,covername)
	local_covername:="./public/"+finalCovername
	fmt.Println(local_filename)
	fmt.Println(local_covername)

	c.SaveUploadedFile(file, local_filename)
	//用第一帧生成cover
	cmd := exec.Command("ffmpeg", "-i", local_filename, "-vframes", "1", "-an",local_covername)
	if cmd.Run() != nil {
		fmt.Println("could not generate frame")
	}

	objectName := model.ObjectName + finalFilename
	fmt.Println(objectName)
	ossCoverName := model.ObjectName + finalCovername
	//上传视频
	err = model.Bucket.PutObjectFromFile(objectName, local_filename)
	if err != nil {
		model.HandleError(err)
	}
	//上传封面
	err = model.Bucket.PutObjectFromFile(ossCoverName, local_covername)
	if err != nil {
		model.HandleError(err)
	}
	//删除本地视频
   err = os.Remove(local_filename)
   if err != nil {
	   fmt.Println("remove local file failed",err)
   }
   //删除本地封面
   err = os.Remove(local_covername)
   if err != nil {
	   fmt.Println("remove local file failed",err)
   }

   title := c.PostForm("title")
   fmt.Println(title)

   video := model.Video_sql{}
//    user_id,_ := c.Get("user_id")
//    uid := config.GetInterfaceToString(user_id)

   video.Author_id= author_id

   play_url := model.Base_url + finalFilename
   coverUrl := model.Base_url + finalCovername
   video.PlayUrl = play_url
   video.CoverUrl = coverUrl
   video.Title = title
   video.CreateTime = time.Now().Unix()

   model.Db_write.Table("videos").Create(&video)
   v:=model.Video_sql{}
   model.Db_write.Table("videos").Where("author_id = ? AND create_time = ? ",video.Author_id,video.CreateTime).First(&v)
   id := strconv.Itoa(int(v.Id))
   //上传视频后，默认给他的点赞数保存半天
   model.Rdb.Set(model.Ctx,id,0,43200*time.Second)

	c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg: "nil",
		})

}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	uid,_ := c.Get("user_id")
	user_id := config.GetInterfaceToString(uid)

	//得到用户的信息
	user := model.User{}
	model.Db_write.Table("user").Where("id = ?",user_id).First(&user)


    video_sql :=[]model.Video_sql{}

	model.Db_write.Table("videos").Where("author_id = ?",user_id).Find(&video_sql)

	video := make([]model.Video,len(video_sql))

	for i := 0 ; i < len(video_sql) ; i++ {

		video[i].Id = video_sql[i].Id

		video[i].Title = video_sql[i].Title
		video[i].CommentCount = video_sql[i].CommentCount
		video[i].FavoriteCount = video_sql[i].FavoriteCount
		video[i].PlayUrl = video_sql[i].PlayUrl
		video[i].CoverUrl = video_sql[i].CoverUrl
		video[i].Author = user
		res:=model.Rdb.SIsMember(model.Ctx,strconv.Itoa(user.Id),strconv.Itoa(int(video[i].Id)))
		isFa,_:=res.Result()
		if isFa{
			video[i].IsFavorite = true
		}else{
			like:=model.Like{}
			model.Db_write.Table("like").Where("video_id = ? AND user_id = ?",video[i].Id,user.Id).First(&like)
			if like.Id!=0{
				video[i].IsFavorite = true
			}else{
				video[i].IsFavorite = false
			}
		}
	}


	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg: "nil",
		},
		VideoList: video,
	})
}
