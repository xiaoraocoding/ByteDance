package middleware

import (
	"fmt"
	"net/http"
	"ByteDance/pkg/app"
	"ByteDance/service"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Query("token")
		fmt.Println(authHeader)
		if authHeader == "" {
			authHeader = c.PostForm("token")
			fmt.Println(authHeader)
		}

		if authHeader == " " {
			c.JSON(http.StatusOK, service.UserResponse{
				Response: service.Response{StatusCode: 1},
			})
			c.Abort()
			return
		}

		mc, err := app.ParseToken(authHeader)
		if err != nil {
			c.JSON(http.StatusOK, service.UserResponse{
				Response: service.Response{StatusCode: 1},
			})
			c.Abort()
			return
		}
		fmt.Println(mc.UserId)
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("user_id", mc.UserId)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息

	}
}
