package task

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	system2 "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// JenkinsJobsWithTest  测试环境 触发构建Jenkins 任务   基于 tasks.json 文件获取任务
func JenkinsJobsWithTest(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, MapName string, taskType string) *modelSystem.JenkinsBuild {
	// 引入不区分大小写的映射
	caseInsensitiveMap := createCaseInsensitiveMap(ExtensionMap)
	// 使用不区分大小写的映射Map
	viewName, _ := getCaseInsensitiveValue(caseInsensitiveMap, MapName)
	global.GVA_LOG.Info("真实视图名: \n" + viewName)
	var tasks map[string][]modelSystem.JenkinsBuild
	jenkinsTasks, err := os.ReadFile("./task/tasks.json")
	if err != nil {
		ReplyWithMessage(bot, webhook, "构建配置文件 tasks.json 未找到！！\n"+err.Error())
		return nil
	}
	// 解析 tasks.json 文件
	err = json.Unmarshal(jenkinsTasks, &tasks)
	if err != nil {
		ReplyWithMessage(bot, webhook, "构建配置文件 tasks.json 解析失败！！\n"+err.Error())
		return nil
	}
	// 输出task 便于调试
	fmt.Printf("解析后的任务配置 %+v\n", tasks)
	// 根据视图名和任务类型查找任务
	for taskViewName, taskList := range tasks {
		// 输出 TaskViewName 便于调试
		global.GVA_LOG.Info("当前的视图名: " + taskViewName)
		// 检查 taskViewName 是否包含输入的 viewName 关键词
		if strings.EqualFold(strings.TrimSpace(taskViewName), strings.TrimSpace(viewName)) {
			fmt.Printf("已匹配的视图名: %s, 任务类型: %s\n", viewName, taskType)
			// 遍历任务列表
			for _, task := range taskList {
				if strings.EqualFold(strings.TrimSpace(task.TaskType), strings.TrimSpace(taskType)) {
					fmt.Printf("输入的任务类型: %s, 任务中的任务类型: %s\n", taskType, task.TaskType)
					fmt.Printf("匹配到任务名: jobName: %v\n", task.JobName)
					// 构建测试环境
					JenkinsBuildJob(bot, webhook, viewName, task.JobName, false) // 构建测试环境
					time.Sleep(30 * time.Second)
					// 调用函数轮询获取任务状态
					go GetJobBuildStatus(bot, webhook, viewName, task.JobName, false)
					return &task
				}
			}
			global.GVA_LOG.Info("未找到对应的 Jenkins Job 任务类型")
			return nil
		}
	}
	global.GVA_LOG.Error("未找到匹配的视图名: ", zap.Any("站点:", viewName))
	ReplyWithMessage(bot, webhook, "测试环境 未找到匹配的站点: "+viewName)

	return nil
}

// ManageService 测试环境 重启应用服务 参数 站点名称 环境参数  命令映射
func ManageService(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, siteName string, EnvName string, command string) {
	type OrderedMap struct {
		Keys   []string
		Values map[string]interface{}
	}
	if siteName == "" {
		log.Printf("站点信息为空")
	}
	//jobName := "exec_command"
	server := "csapi"
	// 获取测试环境执行命令job的参数   false 表示测试环境
	params := GetBuildJobParam(bot, webhook, "", false)
	fmt.Printf("jenkins 参数 parms: %v\n", params)
	// 从数据库中获取站点信息 参数入参  站点名称 或 站点ID   返回值 站点名称对应的一条记录
	var siteConfig modelSystem.YzSiteConfig
	siteConfig.SiteName = siteName
	// 创建 TgService 实例
	tgService := system2.TgService{}
	err := tgService.SelectBySiteName(&siteConfig)
	if err != nil {
		log.Printf("查询数据为空")
	}
	fmt.Printf("对象 %v\n", siteConfig)
	siteID := siteConfig.SiteID
	siteName = siteConfig.SiteName
	siteDomain := siteConfig.Domains
	beginTime := siteConfig.BeginTime
	fmt.Printf("站点ID: %v\n站点名称: %v\n开站时间: %v\n建站时间: %v\n", siteID, siteName, siteDomain, beginTime)
	// 创建表单数据 数据库的值 赋值给 jenkins
	data := url.Values{} // 用于处理查询jenkins URL参数或表单数据
	for key, value := range params {
		data.Set(key.(string), value.(string))
	}
	fmt.Printf("Jenkins赋值给表单data%v\n", data)
	data["site_id"] = []string{strconv.Itoa(int(siteID))}
	data["site"] = []string{siteName}
	data["server"] = []string{server}
	fmt.Printf("jenkins修改后%v\n", data)
}
