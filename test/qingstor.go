package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/qingstor/qingstor-sdk-go/v4/config"
	qs "github.com/qingstor/qingstor-sdk-go/v4/service"
)

func main() {
	// 创建配置
	conf, _ := config.New("MNWBGGEMJFLFQVBENKPH", "jpDwrKzLHINtqEVrK6aafjvjwovhDm7ja8VJlrXy")
	// 初始化青云实例
	qsService, _ := qs.Init(conf)
	bucketName := "bwmobile"
	zoneName := "gd2"
	// 选择存储桶
	bucketService, _ := qsService.Bucket(bucketName, zoneName)
	// 获取桶的权限信息
	if output, err := bucketService.GetACL(); err != nil {
		fmt.Printf("Get acl of bucket(name: %s) failed with given error: %s\n", bucketName, err)
	} else {
		fmt.Printf("The owner of this bucket is %s\n", *output.Owner.ID)
		b, _ := json.Marshal(output.ACL)
		fmt.Println("The acl info of this bucket: ", string(b))
	}
	// 打开要上传的文件
	file, err := os.Open("C:\\Users\\A\\Desktop\\APP\\blgj_107.apk")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	// 获取文件信息 以获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	// 上传文件
	output, err := bucketService.PutObject("/bailigj2/blgj_107.apk", &qs.PutObjectInput{
		Body:          file,
		ContentLength: &fileSize, // 文件大小
	})
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	// 打印上传结果
	fmt.Printf("Upload successful: %v\n", output)
}
