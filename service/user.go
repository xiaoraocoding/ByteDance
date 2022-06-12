package service

import (
	"ByteDance/model"
	"ByteDance/pkg/app"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]model.User{
	"zhangleidouyin": {
		Id: 1,
		//Name:          "zhanglei",
		//FollowCount:   10,
		//FollowerCount: 5,
		//IsFollow:      true,
	},
}

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
// type MyClaims struct {
// 	UserId int `json:"user_id"`
// 	jwt.StandardClaims
// }

// const TokenExpireDuration = time.Hour * 300

// var MySecret = []byte("ThisIsSecret")

// // GenToken 生成JWT
// func GenToken(user_id int) (string, error) {
// 	// 创建一个我们自己的声明
// 	c := MyClaims{
// 		user_id, // 自定义字段
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
// 			Issuer:    "my-project",                               // 签发人
// 		},
// 	}
// 	// 使用指定的签名方法创建签名对象
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
// 	// 使用指定的secret签名并获得完整的编码后的字符串token
// 	return token.SignedString(MySecret)
// }

// // ParseToken 解析JWT
// func ParseToken(tokenString string) (*MyClaims, error) {
// 	// 解析token
// 	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
// 		return MySecret, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
// 		return claims, nil
// 	}
// 	return nil, errors.New("invalid token")
// }

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User model.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := model.User{}

	model.Db_write.Table("user").Where("username = ?", username).Find(&user)
	fmt.Println(user.Id)
	fmt.Println(user.Id)
	fmt.Println(user.Id)

	if user.Id != 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
		})
	} else {
		user = model.User{}
		user.Username = username
		user.Password = password

		model.Db_write.Table("user").Create(&user)
		u := model.User{}
		model.Db_write.Table("user").Where("username = ?", username).First(&u)
		//这里的逻辑是为了后期的点赞等逻辑设计的
		user_follow := strconv.Itoa(u.Id) + "/follow" //当前id的粉丝
		user_action := strconv.Itoa(u.Id) + "/action" //当前id的关注
		fmt.Println(user_action)
		fmt.Println(user_action)

		model.Rdb.Set(model.Ctx, user_follow, 0, -1) //user-follower永久不过期
		model.Rdb.Set(model.Ctx, user_action, 0, -1)

		tokenString, _ := app.GenToken(u.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "nil",
			},
			UserId: int64(u.Id),
			Token:  tokenString,
		})
	}
}

//func Register(c *gin.Context) {
//	username := c.Query("username")
//	password := c.Query("password")
//
//
//
//	token := username + password
//
//	if _, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
//		})
//	} else {
//		atomic.AddInt64(&userIdSequence, 1)
//		newUser := User{
//			Id:   userIdSequence,
//			Name: username,
//		}
//		usersLoginInfo[token] = newUser
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 0},
//			UserId:   userIdSequence,
//			Token:    username + password,
//		})
//	}
//}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := model.User{}

	model.Db_write.Table("user").Where("username = ?", username).Find(&user)
	if user.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "login failed"},
		})
	} else if password == user.Password {
		tokenString, _ := app.GenToken(user.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "nil",
			},
			UserId: int64(user.Id),
			Token:  tokenString,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "login failed"},
		})
	}
}

func UserInfo(c *gin.Context) {

	user_id, _ := c.Get("user_id")
	user := model.User{}
	model.Db_write.Table("user").Where("id = ?", user_id).First(&user)

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "nil",
		},
		User: user,
	})
}

// // JWTAuthMiddleware 基于JWT的认证中间件
// func JWTAuthMiddleware() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
// 		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
// 		// 这里的具体实现方式要依据你的实际业务情况决定
// 		authHeader := c.Query("token")
// 		fmt.Println(authHeader)
// 		if authHeader == "" {
// 			authHeader = c.PostForm("token")
// 			fmt.Println(authHeader)
// 		}

// 		if authHeader == " " {
// 			c.JSON(http.StatusOK, UserResponse{
// 				Response: Response{StatusCode: 1},
// 			})
// 			c.Abort()
// 			return
// 		}

// 		mc, err := app.ParseToken(authHeader)
// 		if err != nil {
// 			c.JSON(http.StatusOK, UserResponse{
// 				Response: Response{StatusCode: 1},
// 			})
// 			c.Abort()
// 			return
// 		}
// 		fmt.Println(mc.UserId)
// 		// 将当前请求的username信息保存到请求的上下文c上
// 		c.Set("user_id", mc.UserId)
// 		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息

// 	}
// }
