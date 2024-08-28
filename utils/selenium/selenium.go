package selenium

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	adminURL               = "https://web.3333c.vip/"
	cookiesUrl             = "https://web.3333c.vip/#/dashboard"
	username               = "yunwei"
	password               = "IRbj2pY27Vm&eMAM"
	captchaFile            = "/Users/david/Downloads/projects/server/captcha.png"
	fullPageScreenshotFile = "/Users/david/Downloads/projects/server/full_page_screenshot.png"
	ApiOCRUrl              = "http://localhost:8000/ocr"
	cookiesFile            = "/Users/david/Downloads/projects/server/cookies.json"
)

func GetAdminLinkTools(siteName string) string {
	var opts []selenium.ServiceOption
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService("/Users/david/tools/chromedriver", 4444, opts...)
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
	time.Sleep(3 * time.Second)
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
		log.Fatalf("截图整个登录页面出错: %v", err)
	}
	// 保存整个页面截图到文件
	if err := os.WriteFile(fullPageScreenshotFile, screenshot, 0644); err != nil {
		log.Fatalf("Error saving the full page screenshot: %v", err)
	}
	// 打开保存的整个页面截图文件
	file, err := os.Open(fullPageScreenshotFile)
	if err != nil {
		log.Fatalf("打开整个登录页面截图报错: %v", err)
	}
	defer file.Close()
	// 解码整个页面截图
	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf("全页面登录页解码报错: %v", err)
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
	fmt.Printf("验证码图像成功保存到文件:%s\n", captchaFile)
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
	time.Sleep(2 * time.Second)
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
		panic(err)
	}
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	data := url.Values{}
	data.Set("image", base64Image)
	data.Set("probability", "false")
	data.Set("png_fix", "false")
	resp, err := http.PostForm(ApiOCRUrl, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("验证码json:  " + string(body))
	// 定义一个结构体或map 来存储解析后的 JSON 数据
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	// 提取验证码对应的data 字段的值
	if data, ok := result["data"].(string); ok {
		fmt.Println("验证码:", data)
		return data, err
	} else {
		fmt.Println("OCR解析验证码错误")
	}
	return "获取OCR验证码错误", nil
}

// GetLinkNoLogin 基于cookie 模式免登录 获取后台链接地址
func GetLinkNoLogin(siteName string) string {
	var opts []selenium.ServiceOption
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService("/Users/david/tools/chromedriver", 4444, opts...)
	if err != nil {
		log.Printf("ChromeDriver server启动错误: %v", err)
	}
	defer service.Stop()
	// 配置 Chrome 浏览器选项以无痕模式启动
	chromeCaps := selenium.Capabilities{
		"browserName": "chrome"}
	chromeOptions := []string{
		"--disable-gpu", // 禁用GPU硬件加速
		//"--headless",    // 启用无界面模式
		"--no-sandbox", // 禁用沙箱命令解决浏览器崩溃问题
		//"--incognito",             // 无痕模式
		"--disable-dev-shm-usage", // 禁用/dev/shm的使用，防止内存共享问题
		"user-data-dir=/Users/david/Library/Application Support/Google/Chrome",
	}
	chromeCaps["goog:chromeOptions"] = map[string]interface{}{
		"args": chromeOptions,
	}
	//  最大化浏览器窗口
	wd, err := selenium.NewRemote(chromeCaps, "http://localhost:4444/wd/hub")
	if err != nil {
		log.Printf("WebDriver 连接失败: %v", err)
	}

	//err = wd.MaximizeWindow("")
	defer wd.Quit()
	// 打开目标网站

	if err := wd.Get(cookiesUrl); err != nil {
		log.Printf("总后台登录首页打开失败: %v"+cookiesUrl, err)
	}

	// 获取所有cookies
	//cookies, _ := wd.GetCookies()
	// 将cookies 序列化 JSON 并保存到文件
	//cookiesJSON, _ := json.Marshal(cookies)
	//os.WriteFile(cookiesFile, cookiesJSON, 0644)
	return "已写入" + siteName
}
