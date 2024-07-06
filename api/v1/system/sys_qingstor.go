package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qingstor/qingstor-sdk-go/v4/config"
	qs "github.com/qingstor/qingstor-sdk-go/v4/service"
)

// 递归获取目录下所有文件
func getAllFiles(bucketService *qs.Bucket, prefix string, allFiles *[]string) error {
	var marker *string

	for {
		listObjectsOutput, err := bucketService.ListObjects(&qs.ListObjectsInput{
			Marker:    marker,
			Prefix:    qs.String(prefix),
			Delimiter: qs.String("/"),
		})
		if err != nil {
			return fmt.Errorf("failed to list objects: %v", err)
		}

		for _, obj := range listObjectsOutput.Keys {
			*allFiles = append(*allFiles, *obj.Key)
		}

		for _, commonPrefix := range listObjectsOutput.CommonPrefixes {
			if err := getAllFiles(bucketService, *commonPrefix, allFiles); err != nil {
				return err
			}
		}

		if listObjectsOutput.NextMarker == nil {
			break
		}
		marker = listObjectsOutput.NextMarker
	}

	return nil
}

func (b *BaseApi) QingCloud(c *gin.Context) {
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
		Prefix:    qs.String(""), // 获取什么前缀开头的目录
	})
	// 获取所有目录
	var dirlists []string
	for _, commonPrefix := range listObjectsOutput.CommonPrefixes {
		dirlists = append(dirlists, *commonPrefix)
		fmt.Println("目录列表:", *commonPrefix)
	}
	// 获取所有站点下apk文件
	//var files []string
	//for _, dir := range dirlists {
	//	files, err := listFilesInDirectory(bucketService, dir)
	//	if err != nil {
	//		log.Fatalf("目录下文件为空%s: %v", dir, err)
	//	}
	//	files = append(files, files...)
	//}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": dirlists,
		//"files": files,
	})
}
