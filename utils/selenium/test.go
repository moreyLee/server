package selenium

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tebeka/selenium"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func GetAdminLinkPhoto(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, siteName string) error {
	var opts []selenium.ServiceOption
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromedriver, 4444, opts...)
	if err != nil {
		log.Printf("ChromeDriver server启动错误: %v", err)
	}
	defer service.Stop()
	// 配置 Chrome 浏览器选项以无痕模式启动
	chromeCaps := selenium.Capabilities{
		"browserName": "chrome"}
	chromeOptions := []string{
		"--disable-gpu", // 禁用GPU硬件加速
		"--headless",    // 开启无界面模式
		"--incognito",   // 无痕模式
		"--windows-size=1920,1080,",
		"--disable-dev-shm-usage", // 禁用/dev/shm的使用，防止内存共享问题
	}
	chromeCaps["goog:chromeOptions"] = map[string]interface{}{
		"args": chromeOptions,
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前执行路径失败: %v", err)
	}
	global.GVA_LOG.Info("当前程序执行路径:  " + dir)
	//  最大化浏览器窗口
	wd, err := selenium.NewRemote(chromeCaps, "http://localhost:4444/wd/hub")
	if err != nil {
		global.GVA_LOG.Info("WebDriver 连接失败: %v", zap.Error(err))
		return err
	}
	defer wd.Quit()

	// 打开管理后台
	if err := wd.Get(adminURL); err != nil {
		global.GVA_LOG.Error("总后台登录首页打开失败:https://web.3333c.vip : %v", zap.Error(err))
		return err
	}
	time.Sleep(3 * time.Second)
	ReplyWithMessage(bot, webhook, siteName+"后台链接正在获取中，请等待30秒...")
	// 获取整个页面截图
	screenshot, err := wd.Screenshot()
	if err != nil {
		global.GVA_LOG.Error("截图整个登录页面出错: %v", zap.Error(err))
		return err
	}
	global.GVA_LOG.Info("截图大小: ", zap.Int("字节", len(screenshot)))
	// 保存整个页面截图到文件
	if err := os.WriteFile(fullPageScreenshotFile, screenshot, 0755); err != nil {
		global.GVA_LOG.Error("保存截图文件失败: %v\n", zap.Error(err))
		return err
	}
	// 打开保存的整个页面截图文件
	file, err := os.Open(fullPageScreenshotFile)
	if err != nil {
		global.GVA_LOG.Error("打开整个登录页面截图报错: %v\n", zap.Error(err))
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			global.GVA_LOG.Error("全局截图文件不存在或解析失败", zap.Error(err))
		}
	}(file)
	photo := tgbotapi.NewPhoto(webhook.Message.Chat.ID, tgbotapi.FilePath(fullPageScreenshotFile))
	_, err = bot.Send(photo)
	if err != nil {
		global.GVA_LOG.Error("发送截图到 Telegram 失败: %v", zap.Error(err))
		return err
	}
	return nil
}
