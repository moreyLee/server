package selenium

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"github.com/tebeka/selenium"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

const (
	adminURL               = "https://web.3333c.vip/"
	username               = "yunwei"
	password               = "IRbj2pY27Vm&eMAM"
	captchaFile            = "D:\\projects\\server\\captcha.png"
	fullPageScreenshotFile = "D:\\projects\\server\\full_page_screenshot.png"
)

func GetAdminLinkTools() {
	var opts []selenium.ServiceOption
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService("D:\\projects\\server\\chromedriver.exe", 4444, opts...)
	if err != nil {
		log.Printf("ChromeDriver server启动错误: %v", err)
		return
	}
	defer service.Stop()
	// 配置 Chrome 浏览器选项以无痕模式启动
	chromeCaps := selenium.Capabilities{
		"browserName": "chrome",
		"goog:chromeOptions": map[string]interface{}{
			"args": []string{
				"--incognito",
			},
		},
	}
	//  最大化浏览器窗口
	wd, err := selenium.NewRemote(chromeCaps, "http://localhost:4444/wd/hub")
	if err != nil {
		log.Printf("WebDriver 连接失败: %v", err)
		return
	}
	err = wd.MaximizeWindow("")
	defer wd.Quit()
	// 打开管理后台
	if err := wd.Get(adminURL); err != nil {
		log.Printf("总后台登录首页打开失败:https://web.3333c.vip : %v", err)
		return
	}
	time.Sleep(3 * time.Second)
	// 刷新页面
	if err := wd.Refresh(); err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
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
		return
	}
	size, err := captchaElem.Size()
	if err != nil {
		fmt.Printf("Error getting captcha size: %v\n", err)
		return
	}
	captchaUrl, _ := captchaElem.GetAttribute("src")
	fmt.Printf("验证码图片链接: %s\n", captchaUrl)
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
	fmt.Printf("成功截图整个登录页面 %s\n", fullPageScreenshotFile)

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
	// 使用 OCR 识别验证码
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(captchaFile)
	captchaText, err := client.Text()
	if err != nil {
		fmt.Printf("Error recognizing captcha: %v\n", err)
		return
	}

	fmt.Printf("Recognized captcha text: %s\n", captchaText)
	//time.Sleep(15 * time.Second)
}
