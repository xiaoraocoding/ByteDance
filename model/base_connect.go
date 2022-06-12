package model

import (
	"ByteDance/conf"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"time"
)

var Db_write *gorm.DB
var err error


func Init() {
	mysql_IP := conf.Get("MYSQL_IP")

	dsn := "root:123456@(" + mysql_IP + ":3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	Db_write, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("write failed", err)
	}



}

var (
	Rdb *redis.Client
)
var Ctx context.Context

// 初始化连接
func InitRedisClient() (err error) {
	Ctx = context.Background()
	addr := conf.Get("redisAddr") + ":" + conf.Get("redisPort")
	fmt.Println(addr)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 100, // 连接池大小
	})

	//这里是设置5秒没链接成功的话，就断开链接
	ctc, cancel := context.WithTimeout(context.Background(),  10*time.Second)
	defer cancel()

	_, err = Rdb.Ping(ctc).Result()
	fmt.Println("start redis success")
	return err
}
