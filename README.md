### 项目背景

本次开发活动由字节跳动官方举办，目标为完成一个轻量级的抖音App



### 安装与使用

- 首先下载本项目的源代码
- 进入项目
- 修改配置文件.env
- sudo go mod tidy



### 使用的开源库

- [gin](https://github.com/gin-gonic/gin) —— 路由、路由组、中间件

- [gorm](https://github.com/go-gorm/gorm) —— ORM 数据操作

- [viper](https://github.com/spf13/viper) —— 配置信息

- [redis](https://github.com/go-redis/redis/v8) —— Redis 操作

- [jwt](https://github.com/dgrijalva/jwt-go) —— JWT 操作

- [oss](https://github.com/aliyun/aliyun-oss-go-sdk/oss)——阿里Oss



### 所有路由

| 请求方法 |         API 地址         | 说明         |
| :------- | :----------------------: | ------------ |
| GET      |       /douyin/feed       | 视频流接口   |
| POST     |  /douyin/user/register/  | 注册         |
| POST     |   /douyin/user/login/    | 登录         |
| POST     |      /douyin/user/       | 用户信息     |
| POST     | /douyin/publish/action/  | 视频投稿接口 |
| GET      |  /douyin/publish/list/   | 视频发布     |
| POST     | /douyin/favorite/action/ | 点赞         |
| GET      |  /douyin/favorite/list/  | 点赞列表     |
| POST     | /douyin/comment/action/  | 评论         |

| GET  |      /douyin/comment/list/      | 评论列表 |
| ---- | :-----------------------------: | -------- |
| POST |    /douyin/relation/action/     | 关注     |
| GET  |  /douyin/relation/follow/list/  | 关注列表 |
| GET  | /douyin/relation/follower/list/ | 粉丝列表 |

