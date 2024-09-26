### http client 库
go get github.com/go-resty/resty/v2
### 删除冲突库
新增CGO_CFLAGS="" 环境变量参数
111ff776c81fe02d30699a285676f3c30c
先删除webhook @CG33333_bot
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/deleteWebhook
修改 telegram webhook  @CG33333_bot hp邮箱
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/setwebhook?url=https://73eb-8-218-67-135.ngrok-free.app/jenkins/telegram-webhook
需要先激活一下测试的webhook  执行一下 telegramwebhook.go

### 打包部署
go build  生成 server 二进制

