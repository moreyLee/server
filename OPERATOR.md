### http client 库
go get github.com/go-resty/resty/v2
### 删除冲突库
新增CGO_CFLAGS="" 环境变量参数
111ff776c81fe02d30699a285676f3c30c
获取TG webhook 状态 @CG33333_bot
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/getWebhookInfo
先删除webhook @CG33333_bot
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/deleteWebhook
修改 telegram webhook  @CG33333_bot hp邮箱
https://api.telegram.org/bot7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4/setwebhook?url=https://a904-217-165-23-20.ngrok-free.app/jenkins/telegram-webhook
需要先激活一下测试的webhook  执行一下 telegramwebhook.go
### 测试机器人 @CG88885_bot
https://api.telegram.org/bot7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4/setWebhook?url=https://orca-awaited-ibex.ngrok-free.app/jenkins/telegram-webhook
https://api.telegram.org/bot7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4/getWebhookInfo
https://api.telegram.org/bot7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4/deleteWebhook

### 打包 
go build -o devops-api
GOOS=linux GOARCH=amd64 go build -o devops-api main.go   # linux 平台
go clean -modcache

## dockerfile 镜像部署
docker  build -t devops-api:v1 . 
# dockerfile 清理构建缓存


# 建立远程分支仓库
git remote add origin http://43.199.1.126:9099/david/devops-api.git
# 查看分支 
git branch 
# 创建并切换分支 
git checkout -b dev 
# 查看远程分支 
git remote -v 
# 记得先切换分支   git push 到公司gitlab
git checkout prod  
#  提交代码  --force 强制推送  覆盖远程仓库   
git push origin prod 
git push origin prod --force 
# 取消缓存中所有文件
git status 

# 查看mysql 活动连接数
SHOW STATUS WHERE `variable_name` = 'Threads_connected';
