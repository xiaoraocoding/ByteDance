package model

import (
	"ByteDance/conf"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

func HandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

var Static_url string
var Base_url string

var ObjectName string
var Bucket *oss.Bucket
func Init_Oss() {
	Endpoint := conf.Get("endpoint")
	AccessKeyId := conf.Get("accessKeyId")
	AccessKeySecret := conf.Get("accessKeySecret")
	BucketName := conf.Get("bucketName")
	ObjectName = conf.Get("objectName")
	Static_url = conf.Get("IMG")
	Base_url = conf.Get("URL")

	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		HandleError(err)
	}

	// 获取存储空间。
	Bucket, err = client.Bucket(BucketName)
	if err != nil {
		HandleError(err)
	}
}
