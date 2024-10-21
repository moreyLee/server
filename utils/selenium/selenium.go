package selenium

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	system2 "github.com/flipped-aurora/gin-vue-admin/server/service/system"
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
	cookiesUrl             = "https://api.3333c.vip/admin/site/config/site/login?siteId=2209" //https://web.3333c.vip/#/site/index" //"https://web.3333c.vip/#/dashboard"
	username               = "yunwei"
	password               = "IRbj2pY27Vm&eMAM"
	captchaFile            = "./code.png"
	fullPageScreenshotFile = "./full_page_screenshot.png"
	ApiOCRUrl              = "http://localhost:8000/ocr"
	chromedriver           = "/Users/david/tools/chromedriver"
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
	time.Sleep(5 * time.Second)
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
	global.GVA_LOG.Info("验证码图像成功保存到文件:  " + captchaFile)
	// 获取OCR 解析的 验证码
	captchaCode, _ := GetCaptchaCode()
	captchaElem, _ = wd.FindElement(selenium.ByID, "captcha")
	captchaElem.SendKeys(captchaCode)
	// 点击登录按钮
	submitElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/main/section/main/div/div/form/div[4]/div/div/span/button")
	submitElem.Click()
	time.Sleep(5 * time.Second)

	//  点击 租户管理
	rentAdminElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/aside[1]/div/div/ul/li[1]/div")
	rentAdminElem.Click()
	time.Sleep(1 * time.Second)
	// 点击 站点管理
	siteAdminElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"/merchants$Menu\"]/li[1]/a")
	siteAdminElem.Click()
	time.Sleep(10 * time.Second)
	// 点击 没有权限 弹窗
	IknowElem, _ := wd.FindElement(selenium.ByCSSSelector, ".ant-modal-confirm-btns")
	IknowElem.Click()
	// 输入 站点名称
	siteNameElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/main/section/main/div/div[3]/div[2]/section/main/div[1]/div[1]/div/span/span/span[2]/input")
	siteNameElem.SendKeys(siteName)
	time.Sleep(1 * time.Second)
	// 点击 搜索按钮
	searchElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/main/section/main/div/div[3]/div[2]/section/main/div[1]/div[7]/button")
	searchElem.Click()
	time.Sleep(1 * time.Second)
	//进入站点  需要加载稍慢
	enterSiteElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/section/main/section/main/div/div[3]/div[2]/section/main/div[3]/div/div/div/div/div/div/div[2]/table/tbody/tr/td[12]/div/div[1]/button[1]")
	enterSiteElem.Click()
	time.Sleep(10 * time.Second)
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
	fmt.Println()
	return siteLink
}

// GetCaptchaCode  获取OCR 解析的验证码
func GetCaptchaCode() (string, error) {
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
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Error("读取OCR 响应失败:  %v", zap.Error(err))
		return "", err
	}
	fmt.Println("验证码json:  " + string(body))
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

// AdminLoginSaveToken 基于token 模式免登录 获取后台链接地址
func AdminLoginSaveToken(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest) string {
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
		"--no-sandbox",
		//"--headless",    // 开启无界面模式
		//"--incognito", // 无痕模式
		//"--windows-size=1920,1080,",
		"--disable-dev-shm-usage", // 禁用/dev/shm的使用，防止内存共享问题
	}
	chromeCaps["goog:chromeOptions"] = map[string]interface{}{
		"args": chromeOptions,
	}
	//  最大化浏览器窗口
	wd, err := selenium.NewRemote(chromeCaps, "http://localhost:4444/wd/hub")
	if err != nil {
		log.Printf("WebDriver 连接失败: %v", err)
	}

	err = wd.MaximizeWindow("")
	defer wd.Quit()

	// 打开管理后台
	if err := wd.Get(adminURL); err != nil {
		log.Printf("总后台登录首页打开失败:https://web.3333c.vip : %v", err)
	}
	// 刷新页面
	if err := wd.Refresh(); err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	// 输入用户名
	usernameElem, _ := wd.FindElement(selenium.ByID, "username")
	usernameElem.SendKeys(username)
	// 输入密码
	passwordElem, _ := wd.FindElement(selenium.ByID, "password")
	passwordElem.SendKeys(password)
	// 等待60秒 手动输入验证码
	time.Sleep(10 * time.Second)
	// 登录
	submitElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"root\"]/section/main/section/main/div/div/form/div[4]/div/div/span/button")
	submitElem.Click()
	time.Sleep(15 * time.Second)
	// 登录成功后  检查当前url  确认是否成功
	currentURL, _ := wd.CurrentURL()
	fmt.Printf("当前页面URL:%s\n", currentURL)
	// 使用 js 获取 local Storage 中的值
	token, _ := wd.ExecuteScript("return window.localStorage.getItem('token');", nil)
	fmt.Printf("token值%s\n", token)
	//
	// 将登录后的token 保存到数据库
	tgService := system2.TgService{}
	err = tgService.SaveAdminLoginToken(token.(string))
	if err != nil {
		ReplyWithMessage(bot, webhook, "写入数据库报错:\n"+err.Error())
	}
	// 打开一个新的标签页 并导航到新地址
	newTitle := "window.open('https://web.3333c.vip/#/site/index', '_blank');"
	wd.ExecuteScript(newTitle, nil)
	// 切换到新打开的标签页
	handles, _ := wd.WindowHandles()
	// 切换到新打开的标签页
	if len(handles) > 1 {
		wd.SwitchWindow(handles[len(handles)-1])
	}
	time.Sleep(60 * time.Second)
	return token.(string)
}

// GetLinkNoLogin 免登录 获取后台链接地址
func GetLinkNoLogin() string {
	var opts []selenium.ServiceOption
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromedriver, 4444, opts...)
	if err != nil {
		log.Printf("ChromeDriver server启动错误: %v", err)
	}
	defer service.Stop()
	// 配置 Chrome 浏览器选项
	chromeCaps := selenium.Capabilities{
		"browserName": "chrome"}
	chromeOptions := []string{
		"--disable-gpu", // 禁用GPU硬件加速
		"--no-sandbox",
		//"--headless",    // 开启无界面模式
		//"--incognito", // 无痕模式
		//"--windows-size=1920,1080,",
		"--disable-dev-shm-usage", // 禁用/dev/shm的使用，防止内存共享问题
		//"--user-data-dir=/Users/david/Library/Application Support/Google/Chrome",
	}

	chromeCaps["goog:chromeOptions"] = map[string]interface{}{
		"args": chromeOptions,
	}
	//  最大化浏览器窗口
	wd, _ := selenium.NewRemote(chromeCaps, "http://localhost:4444/wd/hub")
	err = wd.MaximizeWindow("")
	defer wd.Quit()
	// 打开目标网页
	wd.Get("https://web.3333c.vip/#/site/index")
	time.Sleep(5 * time.Second)
	wd.Refresh()

	// 创建一个新的cookie
	tokenstr := "ab60b1a668c8caeeef41f25691a79694"
	cookie := &selenium.Cookie{Name: "HTTP_TOKEN", Value: tokenstr, Path: "/"}

	// 将token 添加到cookies 中
	if err := wd.AddCookie(cookie); err != nil {
		panic(err)
	}
	// 刷新页面 让cookies 生效
	time.Sleep(10 * time.Second)
	if err := wd.Refresh(); err != nil {
		panic(err)
	}
	// 查询数据库获取登录token
	tgService := system2.TgService{}
	token, err := tgService.GetAdminLoginToken() // 返回的类型为对象
	fmt.Printf("获取数据库的token%v 报错信息:%v\n", token, err)
	fmt.Printf("数据库httpToken值: %v\n", token.HttpToken)
	// 使用 js 获取 local Storage 中的值
	//nologintoken, _ := wd.ExecuteScript("return window.localStorage.getItem('token');", nil)
	//fmt.Printf("登录前的token值%s\n", nologintoken)
	// 将 Token 写入Local Storage
	//script := fmt.Sprintf("window.localStorage.setItem('token', '%s');", token.HttpToken)
	//fmt.Println("script执行情况:  ", script)
	//_, err = wd.ExecuteScript(script, nil)
	//if err != nil {
	//	fmt.Println("设置local Storage 报错", err)
	//}
	//tokenstr, ok := token.HttpToken.(string)
	//_, err = wd.ExecuteScript(`window.localStorage.setItem('token', arguments[0]);`, []interface{}{token.HttpToken})
	//fmt.Printf("写入浏览器的token%v\n", []interface{}{token.HttpToken})
	//if err != nil {
	//	log.Fatalf("设置token 到localStorage 报错:%v", err)
	//}
	// 打开一个新的标签页 并导航到新地址
	//newTitle := "window.open('https://web.3333c.vip/#/site/index', '_blank');"
	//wd.ExecuteScript(newTitle, nil)
	// 切换到新打开的标签页
	//handles, _ := wd.WindowHandles()
	// 切换到新打开的标签页
	//if len(handles) > 1 {
	//	wd.SwitchWindow(handles[len(handles)-1])
	//}
	//time.Sleep(15 * time.Second)
	// 刷新页面 已使token生效 验证免登录
	//wd.Refresh()
	time.Sleep(900 * time.Second)
	// 访问目标页面 验证免登录
	//wd.Get("https://web.3333c.vip/#/site/index")
	//time.Sleep(3600 * time.Second)
	return token.HttpToken
}
