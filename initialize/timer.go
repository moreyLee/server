package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func Timer() {
	//启动协程运行函数
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 其他定时任务定在这里 参考上方使用方法
		// 定时响应 @telegram 机器人的消息 调用三方接口
		//_, err = global.GVA_Timer.AddTaskByFunc("定时任务标识", "@every 1h", func() {
		//	task.SendMessage()
		//	//if err != nil {
		//	//	fmt.Println("定时任务机器人错误", err)
		//	//}
		//}, "@机器人响应并执行三方接口", option...)
		//if err != nil {
		//	fmt.Println("添加定时任务错误:", err)
		//}
	}()
}
