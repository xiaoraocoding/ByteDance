package service

import (
	"ByteDance/config"
	"ByteDance/mes"
	"ByteDance/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []model.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {

	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	to_uid, _ := strconv.Atoi(to_user_id)

	action_t, _ := strconv.Atoi(action_type)

	uid, _ := c.Get("user_id")

	user_id := config.GetInterfaceToString(uid)
	user_i, _ := strconv.Atoi(user_id)
	var follow model.Follow
	if action_t == 1 { //关注
		subscribe := model.Subscribe{}
		subscribe.BeSubscribe = int64(to_uid) //被关注的id
		subscribe.Subscribe = int64(user_i)   //关注者的id
		subscribe.IsDel = 0                   //说明是关注
		timeStr := time.Now().Format("2006-01-02")
		subscribe.CreatTime = timeStr
		// model.Db_write.Table("db_subscribe").Create(&subscribe)
		follow_sql := model.Follow_sql{}
		//表中并没有设置userid主键，防止搞鬼，让数据库出现多条相同数据
		model.Db_write.Table("follow").Where("user_id = ? AND target_id =?", user_id, follow.ToUserId).First(&follow_sql)
		if follow_sql.Id == 0 {
			u_id, _ := strconv.Atoi(user_id)
			follow_sql.UserId = int64(u_id)
			follow_sql.TargetId = follow.ToUserId
			model.Db_write.Table("follow").Create(&follow_sql)
			followed_sql := model.Followed_sql{
				UserId:   follow.ToUserId, //交换了，userId在followed表中成为userId
				TargetId: int64(u_id),
			}
			user_follow := to_user_id + "/follow" //粉丝
			user_action := user_id + "/action" //当前用户的关注数量
	
			model.Rdb.Incr(model.Ctx, user_follow) //当前用户的粉丝加1
			value := model.Rdb.Get(model.Ctx, user_follow)
			model.Rdb.Incr(model.Ctx, user_action)
			value_action := model.Rdb.Get(model.Ctx, user_action)
	
			value_result, _ := value.Result()
			r_value, _ := strconv.Atoi(value_result)
			value_action_result, _ := value_action.Result()
			r_action_value, _ := strconv.Atoi(value_action_result)
			message:=mes.Message{
				Followed:followed_sql,
				FollowCount: r_action_value,
				FollowerCount: r_value,

			}
			user_list := to_user_id + "/follower/list"   //粉丝的列表
			user_follow_list := user_id + "/follow/list" //当前用户的关注列表
			model.Rdb.SAdd(model.Ctx, user_list, user_i)
			model.Rdb.SAdd(model.Ctx, user_follow_list, to_uid)
			if err := mes.Producer(message, "follow_back"); err != nil {
				return
			}
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "",
			})
		}
	} else if action_t == 2 {
		// subscribe := model.Subscribe{}
		// model.Db_write.Table("db_subscribe").Model(&subscribe).Where("be_subscribe = ? and subscribe = ?",int64(to_uid),int64(user_i)).Update("is_del",1)
		user_follow := to_user_id + "/follow"
		user_action := user_id + "/action"

		model.Rdb.IncrBy(model.Ctx, user_follow, -1) //当前用户的粉丝减1
		value := model.Rdb.Get(model.Ctx, user_follow)
		model.Rdb.IncrBy(model.Ctx, user_action, -1)
		value_action := model.Rdb.Get(model.Ctx, user_action)

		value_result, _ := value.Result()
		r_value, _ := strconv.Atoi(value_result)
		value_action_result, _ := value_action.Result()
		r_action_value, _ := strconv.Atoi(value_action_result)


		follow_sql := model.Follow_sql{}
		u_id, _ := strconv.Atoi(user_id)
		follow_sql.UserId = int64(u_id)
		follow_sql.TargetId = follow.ToUserId
		model.Db_write.Table("follow").Where("user_id = ? AND target_id = ?", u_id, follow.ToUserId).Delete(&follow_sql)
		followed_sql := model.Followed_sql{
			UserId:   follow.ToUserId,
			TargetId: int64(u_id),
		}
		message:=mes.Message{
			Followed:followed_sql,
			FollowCount: r_action_value,
			FollowerCount: r_value,

		}
		// user := model.User{}
		// model.Db_write.Model(&user).Table("user").Where("id = ?", to_uid).Update("follower_count", r_value)
		// us := model.User{}
		// model.Db_write.Model(&us).Table("user").Where("id = ?", user_i).Update("follow_count", r_action_value)


		user_list := to_user_id + "/follower/list"   //粉丝的列表
		user_follow_list := user_id + "/follow/list" //当前用户的关注列表
		model.Rdb.SRem(model.Ctx, user_list, user_i)
		model.Rdb.SRem(model.Ctx, user_follow_list, to_uid)
		if err := mes.Producer(message, "cancer_back"); err != nil {
			println(err.Error())
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "",
		})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {

	uid, _ := c.Get("user_id")


	user_i := config.GetInterfaceToString(uid)

	user_id := c.Query("user_id")
	fmt.Println(user_id)
	uid_follow_list := user_i + "/follow/list"
	user_follow_list := user_id + "/follow/list" //关注列表

	result := model.Rdb.SMembers(model.Ctx,user_follow_list)
	res_list := result.Val()

	result_list := []int{}

	for i := 0 ; i < len(res_list) ; i ++ {
		res,_ := strconv.Atoi(res_list[i])
		result_list = append(result_list,res)
	}


	var userList []model.User

	for i := 0 ; i < len(result_list) ; i ++ {
		user := model.User{}
		//model.Db_write.Table("user").Model(&user).Where("id = ?",result_list[i])
		model.Db_write.Table("user").Where("id = ?",result_list[i]).Find(&user)
		boolRes := model.Rdb.SIsMember(model.Ctx,uid_follow_list,result_list[i])
		boolres := boolRes.Val()
		user.IsFollow = boolres
		userList = append(userList,user)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")

	user_id_follower_list := user_id + "/follower/list"

	uid, _ := c.Get("user_id")
	user_i := config.GetInterfaceToString(uid)
	uid_follow_list := user_i + "/follow/list"

	result := model.Rdb.SMembers(model.Ctx,user_id_follower_list)
	res_list := result.Val()

	result_list := []int{}

	for i := 0 ; i < len(res_list) ; i ++ {
		res,_ := strconv.Atoi(res_list[i])
		result_list = append(result_list,res)
	}

	var userList []model.User

	for i := 0 ; i < len(result_list) ; i ++ {
		user := model.User{}
		//model.Db_write.Table("user").Model(&user).Where("id = ?",result_list[i])
		model.Db_write.Table("user").Where("id = ?",result_list[i]).Find(&user)
		boolRes := model.Rdb.SIsMember(model.Ctx,uid_follow_list,result_list[i])
		boolres := boolRes.Val()
		user.IsFollow = boolres
		userList = append(userList,user)
	}


	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
