package douyin

import (
	"ByteDance/middleware"
	"ByteDance/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", middleware.JWTAuthMiddleware(), service.Feed)
	apiRouter.GET("/user/", middleware.JWTAuthMiddleware(), service.UserInfo)
	apiRouter.POST("/user/register/", service.Register)
	apiRouter.POST("/user/login/", service.Login)
	apiRouter.POST("/publish/action/", middleware.JWTAuthMiddleware(), service.Publish)
	apiRouter.GET("/publish/list/", middleware.JWTAuthMiddleware(), service.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.JWTAuthMiddleware(), service.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.JWTAuthMiddleware(), service.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.JWTAuthMiddleware(), service.CommentAction)
	apiRouter.GET("/comment/list/", middleware.JWTAuthMiddleware(), service.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.JWTAuthMiddleware(), service.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.JWTAuthMiddleware(), service.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.JWTAuthMiddleware(), service.FollowerList)

}
