### http client 库
go get github.com/go-resty/resty/v2
### 删除冲突库
新增CGO_CFLAGS="" 环境变量参数
111ff776c81fe02d30699a285676f3c30c
获取TG webhook 状态
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/getWebhookInfo
先删除webhook @CG33333_bot
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/deleteWebhook
修改 telegram webhook  @CG33333_bot hp邮箱
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/setwebhook?url=https://a904-217-165-23-20.ngrok-free.app/jenkins/telegram-webhook
需要先激活一下测试的webhook  执行一下 telegramwebhook.go
### 测试机器人 @CG88885_bot
https://api.telegram.org/bot7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4/setWebhook?url=https://e498-217-165-23-20.ngrok-free.app/telegram-webhook
https://api.telegram.org/bot7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4/getWebhookInfo
https://api.telegram.org/bot7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4/deleteWebhook
### 打包部署
go build  生成 server 二进制
### 打包成指定的包名称
go build -o devops-api
## 镜像部署 基于dockerfile 

