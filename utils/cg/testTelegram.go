package cg

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

const (
	telegramBotToken = "7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4"
	jenkinsURL       = "http://192.168.217.128:8082"
	jenkinsUser      = "admin"
	jenkinsAPIToken  = "11c9bc0d6ea88891f45ee4cfe5bd218287"
)

func main() {
	// 创建新的机器人实例
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	// 设置调试模式
	bot.Debug = true

	// 获取更新
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// 处理更新
	for update := range updates {
		if update.Message != nil { // 如果有消息更新
			// 检查消息是否是命令
			if update.Message.IsCommand() {
				handleCommand(bot, update.Message)
			}
		}
	}
}

// 处理命令
func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// 获取命令及参数
	command := message.Command()
	args := message.CommandArguments()

	switch command {
	case "build":
		// 处理 /build 命令
		if args == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Usage: /build <job_name>")
			bot.Send(msg)
			return
		}
		jobName := args
		err := triggerJenkinsJob(jobName)
		var response string
		if err != nil {
			response = fmt.Sprintf("Failed to trigger Jenkins job '%s': %v", jobName, err)
		} else {
			response = fmt.Sprintf("Successfully triggered Jenkins job '%s'", jobName)
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
	default:
		// 未知命令
		msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command")
		bot.Send(msg)
	}
}

// 触发 Jenkins 构建任务
func triggerJenkinsJob(jobName string) error {
	client := resty.New()
	client.SetBasicAuth(jenkinsUser, jenkinsAPIToken)

	url := fmt.Sprintf("%s/job/%s/build", jenkinsURL, jobName)
	resp, err := client.R().Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
