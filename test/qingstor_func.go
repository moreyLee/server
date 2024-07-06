package main

import (
	"fmt"
	"github.com/qingstor/qingstor-sdk-go/v4/config"
	qs "github.com/qingstor/qingstor-sdk-go/v4/service"
)

func main() {
	conf, _ := config.New("MNWBGGEMJFLFQVBENKPH", "jpDwrKzLHINtqEVrK6aafjvjwovhDm7ja8VJlrXy")
	// 初始化青云实例
	qsService, _ := qs.Init(conf)
	bucketName := "bwmobile"
	zoneName := "gd2"
	// 选择存储桶
	bucketService, _ := qsService.Bucket(bucketName, zoneName)
	// 获取APP目录列表
	listObjectsOutput, _ := bucketService.ListObjects(&qs.ListObjectsInput{
		Delimiter: qs.String("/"),
		//Prefix:    qs.String("app/"),
	})
	fmt.Println("对象存储根目录文件列表:")
	// 打印对象列表
	for _, object := range listObjectsOutput.Keys {
		fmt.Println("文件列表:", *object.Key)
	}
	// 获取所有目录
	for _, commonPrefix := range listObjectsOutput.CommonPrefixes {
		fmt.Println("目录列表:", *commonPrefix)
	}

}
