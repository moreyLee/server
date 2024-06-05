package main

import (
	"fmt"
	"github.com/qingstor/qingstor-sdk-go/v4/config"
	"github.com/qingstor/qingstor-sdk-go/v4/service"
	"os"
	"path/filepath"
)

// 定义全局变量
var (
	accessKey    = os.Getenv("MNWBGGEMJFLFQVBENKPH")
	secretKey    = os.Getenv("jpDwrKzLHINtqEVrK6aafjvjwovhDm7ja8VJlrXy")
	bucketName   = "bwmobile"
	zoneName     = "gd2"
	objectKey    string
	filePath     string
	baseAddress  = "https://"
	domainSuffix = ".qingstor.com/"
)

// GetBucketService 函数用于获取指定的 Bucket 服务
func GetBucketService(qsService *service.Service, bucketName, zoneName string) (*service.Bucket, error) {
	// 获取指定桶的 Bucket 对象
	bucketService, err := qsService.Bucket(bucketName, zoneName)
	if err != nil {
		fmt.Println("Failed to get bucket service:", err)
		return nil, err
	}

	return bucketService, err
}

// MatchFiles 函数用于匹配指定文件
func MatchFiles(filePattern string) ([]string, error) {
	files, err := filepath.Glob(filePattern)
	if err != nil {
		fmt.Println("Failed to match files:", err)
		return nil, fmt.Errorf("failed to match files: %v", err)
	}
	return files, nil
}

// UploadFile 函数用于上传文件
func UploadFile(bucketService *service.Bucket, objectKey, filePath string) ([]string, error) {
	// 用于存储上传成功的文件列表
	var uploadedFiles []string

	// 匹配文件列表函数
	matches, err := MatchFiles(filePath)
	if err != nil {
		return nil, err
	}

	// 遍历匹配的文件列表，逐个上传
	for _, match := range matches {
		// 打开文件
		file, err := os.Open(match)
		if err != nil {
			return uploadedFiles, err
		}

		// 延迟关闭文件
		defer file.Close()

		// 获取文件信息
		fileInfo, err := file.Stat()
		if err != nil {
			return uploadedFiles, err
		}

		// 上传对象文件
		output, err := bucketService.PutObject(objectKey+"/"+fileInfo.Name(), &service.PutObjectInput{Body: file})
		if err != nil {
			return uploadedFiles, err
		}

		// 所有 >= 400 的状态码都会被认为是错误
		if *output.StatusCode >= 400 {
			return uploadedFiles, fmt.Errorf("Failed to upload file: %d", fileInfo.Name())
		}

		// 上传成功，将文件名添加到上传成功的文件列表中
		uploadedFiles = append(uploadedFiles, fileInfo.Name())
	}

	return uploadedFiles, nil
}

// InitConfig 函数用于初始化配置
func InitConfig() (*service.Service, error) {
	// 创建配置
	configuration, err := config.New(accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	fmt.Println(configuration)

	// 初始化服务
	qsService, err := service.Init(configuration)
	if err != nil {
		fmt.Println("Failed to initialize QingStor service:", err)
		return nil, err
	}
	return qsService, nil

}

func main() {
	// 初始化配置
	qingStor, err := InitConfig()
	if err != nil {
		fmt.Println("Failed to create configuration:", err)
		return
	}

	// 创建对象服务
	bucketService, err := GetBucketService(qingStor, bucketName, zoneName)
	if err != nil {
		fmt.Println("Failed to create bucket service:", err)
		return
	}

	// 调用上传函数
	uploadedFiles, err := UploadFile(bucketService, objectKey, filePath)
	if err != nil {
		fmt.Println("Failed to upload files:", err)
		return
	}

	// 打印上传成功的文件信息
	fmt.Println("Uploaded files:")
	for _, file := range uploadedFiles {
		fmt.Println(file)
	}
}
