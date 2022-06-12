package main

import (
	"ByteDance/conf"
	"ByteDance/mes"
	"flag"
	"fmt"

	"ByteDance/config"
	"ByteDance/model"
	"ByteDance/router/douyin"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.Init()
	var env string
	flag.StringVar(&env, "env", "","自己配置 如--env=test")
	//如果是 --env=test的话那么就会读取.env.test文件
	flag.Parse()
	conf.InitConfig(env)
	mes.InitConsumer("follow_back","test")
	mes.InitCancerConsumer("cancer_back","canceltest")

	douyin.InitRouter(r)
	model.Init()
	err := model.InitRedisClient()
	if err != nil {
		fmt.Println("init redis",err)
	}

	model.Init_Oss()
	port := conf.Get("app.port")

	r.Run(":" + port)
}
