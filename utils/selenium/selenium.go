package selenium

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tebeka/selenium"
	"go.uber.org/zap"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"
)

const (
	adminURL               = "https://web.3333c.vip/"
	username               = "yunwei"
	password               = "IRbj2pY27Vm&eMAM"
	captchaFile            = "./code.png"
	fullPageScreenshotFile = "./full_page_screenshot.png"
	ApiOCRUrl              = "http://localhost:8000/ocr"
	chromedriver           = "/Users/david/tools/chromedriver/chromedriver"
	port                   = 5555
)

// ReplyWithMessage 全局引用 用于小飞机发送消息
func ReplyWithMessage(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, message string) {
	replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
	replyText.ReplyToMessageID = webhook.Message.MessageID
	_, _ = bot.Send(replyText)
}

// GetAdminLinkTools 登录后获取链接地址
func GetAdminLinkTools(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, siteName string) string {
	var opts []selenium.ServiceOption
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromedriver, port, opts...)
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
	wd, err := selenium.NewRemote(chromeCaps, "http://localhost:5555/wd/hub")
	if err != nil {
		global.GVA_LOG.Error("WebDriver 连接失败: %v", zap.Error(err))
		errorMessage := fmt.Sprintf("WebDriver 连接失败: %v", zap.Error(err))
		ReplyWithMessage(bot, webhook, errorMessage)
	}
	defer wd.Quit()

	// 打开管理后台
	if err := wd.Get(adminURL); err != nil {
		global.GVA_LOG.Info("总后台登录首页打开失败:https://web.3333c.vip : %v", zap.Error(err))
	}
	time.Sleep(3 * time.Second)
	ReplyWithMessage(bot, webhook, siteName+"后台链接正在获取中，请等待一分钟...")
	// 刷新页面
	if err := wd.Refresh(); err != nil {
		global.GVA_LOG.Error("刷新页面报错", zap.Error(err))
		return ""
	}
	time.Sleep(3 * time.Second)
	// 输入用户名
	usernameElem, _ := wd.FindElement(selenium.ByID, "username")
	usernameElem.SendKeys(username)
	// 输入密码
	passwordElem, _ := wd.FindElement(selenium.ByID, "password")
	passwordElem.SendKeys(password)
	// 查找验证码图像元素 基于css 样式选择器 样式唯一
	captchaElem, err := wd.FindElement(selenium.ByCSSSelector, ".login-captcha")
	if err != nil {
		panic(err)
	}
	// 获取验证码图像的位置和大小
	loc, err := captchaElem.Location()
	if err != nil {
		fmt.Printf("Error getting captcha location: %v\n", err)
	}
	size, err := captchaElem.Size()
	if err != nil {
		fmt.Printf("Error getting captcha size: %v\n", err)
	}
	fmt.Printf("截取验证码的坐标位置:%v\n,\n", loc)
	// 获取整个页面截图
	screenshot, err := wd.Screenshot()
	if err != nil {
		global.GVA_LOG.Error("截图整个登录页面出错: %v", zap.Error(err))
		return ""
	}
	global.GVA_LOG.Info("截图大小: ", zap.Int("字节", len(screenshot)))
	// 保存整个页面截图到文件
	if err := os.WriteFile(fullPageScreenshotFile, screenshot, 0755); err != nil {
		global.GVA_LOG.Error("保存截图文件失败: %v\n", zap.Error(err))
	}
	// 打开保存的整个页面截图文件
	file, err := os.Open(fullPageScreenshotFile)
	if err != nil {
		global.GVA_LOG.Error("打开整个登录页面截图报错: %v\n", zap.Error(err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			global.GVA_LOG.Error("全局截图文件不存在或解析失败", zap.Error(err))
		}
	}(file)
	// 解码整个页面截图
	img, err := png.Decode(file)
	if err != nil {
		global.GVA_LOG.Error("全页面登录页解码报错: %v\n", zap.Error(err))
	}
	// 裁剪出验证码区域
	bounds := image.Rect(loc.X, loc.Y, loc.X+size.Width, loc.Y+size.Height)
	captchaImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(bounds)
	// 保存裁剪后的验证码图像到文件
	captchaOutputFile, err := os.Create(captchaFile)
	if err != nil {
		log.Fatalf("保存截取后的验证码图像到文件报错Error: %v", err)
	}
	defer captchaOutputFile.Close()
	if err := png.Encode(captchaOutputFile, captchaImg); err != nil {
		log.Fatalf("验证码图像解码报错Error: %v", err)
	}
	global.GVA_LOG.Info("验证码图像成功保存到文件:  ")
	// 获取OCR 解析的 验证码
	captchaCode, err := GetCaptchaCode(bot, webhook)
	if err != nil {
		global.GVA_LOG.Error("OCR解析验证码失败: \n ", zap.Error(err))
		ReplyWithMessage(bot, webhook, fmt.Sprintf("OCR验证码验证失败，请稍后再试: %v", err))
		return ""
	}
	captchaElem, err = wd.FindElement(selenium.ByID, "captcha")
	if err != nil {
		global.GVA_LOG.Error("找不到验证码输入框:", zap.Error(err))
		ReplyWithMessage(bot, webhook, "找不到验证码输入框，请刷新页面")
		return ""
	}

	captchaElem.SendKeys(captchaCode)
	// 点击登录按钮
	submitElem, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/main/section/main/div/div/form/div[4]/div/div/span/button")
	submitElem.Click()
	if err != nil {
		code, _ := GetCaptchaCode(bot, webhook)
		ReplyWithMessage(bot, webhook, "登录失败 验证识别错误:"+code+"\n"+err.Error())
		return ""
	}
	time.Sleep(3 * time.Second) // 登录完成

	//弹窗 商户首页数据异常
	//excepElem, err := wd.FindElement(selenium.ByCSSSelector, ".ant-modal-confirm-btns")
	//excepElem.Click()

	//  点击 租户管理
	rentAdminElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/aside[1]/div/div/ul/li[1]/div")
	rentAdminElem.Click()
	time.Sleep(1 * time.Second)
	global.GVA_LOG.Info(fmt.Sprintf("点击租户管理 %s\n 等待1秒", rentAdminElem))

	// 点击 站点管理
	siteAdminElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"/merchants$Menu\"]/li[1]/a")
	siteAdminElem.Click()
	time.Sleep(10 * time.Second)
	global.GVA_LOG.Info(fmt.Sprintf("点击站点管理 %s\n 等待10秒", siteAdminElem))

	// 点击 没有权限 弹窗
	IknowElem, _ := wd.FindElement(selenium.ByCSSSelector, ".ant-modal-confirm-btns")
	IknowElem.Click()
	global.GVA_LOG.Info(fmt.Sprintf("点击没权限弹窗 %s\n ", IknowElem))

	// 输入 站点名称
	siteNameElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/main/section/main/div/div[3]/div[2]/section/main/div[1]/div[1]/div/span/span/span[2]/input")
	siteNameElem.SendKeys(siteName)
	time.Sleep(1 * time.Second)
	global.GVA_LOG.Info(fmt.Sprintf("点击 站点名称 等待1秒 %s\n ", siteNameElem))

	// 点击 搜索按钮
	searchElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/main/section/main/div/div[3]/div[2]/section/main/div[1]/div[7]/button")
	searchElem.Click()
	time.Sleep(1 * time.Second)
	global.GVA_LOG.Info(fmt.Sprintf("点击 搜索按钮 等待1秒 %s\n ", searchElem))

	//进入站点  需要加载稍慢
	enterSiteElem, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/main/section/main/div/div[3]/div[2]/section/main/div[3]/div/div/div/div/div/div/div[2]/table/tbody/tr/td[12]/div/div[1]/button[1]")
	if err != nil {
		ReplyWithMessage(bot, webhook, "站名名不存在: "+siteName+",请检查站点名称")
		return ""
	}
	enterSiteElem.Click()
	time.Sleep(10 * time.Second)
	global.GVA_LOG.Info(fmt.Sprintf("点击 进入站点 等待10秒 %s\n ", enterSiteElem))

	// 获取站点地址链接
	handles, _ := wd.WindowHandles() // 获取当前所有窗口句柄
	// 点击打开一个新窗口 切换到新窗口
	if len(handles) > 1 {
		newWindowHandle := handles[len(handles)-1]
		wd.SwitchWindow(newWindowHandle)
	}
	// 获取 新打开页面的站点url
	siteLink, _ := wd.CurrentURL()
	fmt.Println(siteName+"站点地址: ", siteLink)
	ReplyWithMessage(bot, webhook, siteName+"站点地址:\n"+siteLink)
	return siteLink
}

// GetCaptchaCode  获取OCR 解析的验证码
func GetCaptchaCode(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest) (string, error) {
	imageData, err := os.ReadFile(captchaFile)
	if err != nil {
		global.GVA_LOG.Error("读取验证码文件失败: %v", zap.Error(err))
		return "", err
	}
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	data := url.Values{}
	data.Set("image", base64Image)
	data.Set("probability", "false")
	data.Set("png_fix", "false")
	resp, err := http.PostForm(ApiOCRUrl, data)
	if err != nil {
		global.GVA_LOG.Error("发送OCR 请求失败: %v", zap.Error(err))
		ReplyWithMessage(bot, webhook, "OCR 8000 端口请求失败，请检查服务是否正常！\n"+err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Error("读取OCR 响应失败:  %v", zap.Error(err))
		return "", err
	}
	//fmt.Println("验证码json:  " + string(body))
	// 定义一个结构体或map 来存储解析后的 JSON 数据
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		global.GVA_LOG.Error("解析OCR 响应的 JSON 数据失败:  %v", zap.Error(err))
		return "", err
	}
	// 提取验证码对应的data 字段的值
	if data, ok := result["data"].(string); ok {
		global.GVA_LOG.Info("\n 原始验证码: " + data)
		// 过滤非数字字符，只保留数字
		onlyDigits := strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) {
				return r
			}
			return -1
		}, data)
		global.GVA_LOG.Info("\n 过滤后的验证码:  " + onlyDigits)
		return onlyDigits, err
		//
	} else {
		global.GVA_LOG.Info("OCR 解析验证码失败: %v", zap.Any("response", result))
		return "获取OCR验证码错误", nil
	}
}
