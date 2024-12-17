package cron

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/robfig/cron/v3"
)

// RegisterCron 注册定时任务
func RegisterCron(c *cron.Cron) {
	// 注册 SyncJenkinsData 定时任务，每小时执行一次
	_, err := c.AddFunc("@hourly", func() {
		task.SyncJenkinsData() // 调用定时任务
	})
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("Failed to register SyncJenkinsData task: %v", err))
		return
	}

	global.GVA_LOG.Info("Cron tasks registered successfully!")
}
