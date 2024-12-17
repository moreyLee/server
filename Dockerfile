FROM golang:alpine as builder

WORKDIR /app
# 复制源代码并编译
# 在 Dockerfile 顶部添加
ARG BUILD_TIMESTAMP
ARG BUILD_VERSION

# 作为环境变量写入镜像
ENV BUILD_TIMESTAMP=$BUILD_TIMESTAMP
ENV BUILD_VERSION=$BUILD_VERSION

COPY . .
# 程序编译打包 设置环境变量
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://proxy.golang.org,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && echo "Build Timestamp: $BUILD_TIMESTAMP" > .build_info \
    && echo "Build Version: $BUILD_VERSION" >> .build_info \
    && go build -o devops-api .

FROM alpine:latest

LABEL MAINTAINER="David588@gmail.com"

WORKDIR /app
# 安装 tzdata 以支持时区配置
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Dubai /etc/localtime \
    && echo "Asia/Dubai" > /etc/timezone
# 测试环境 Jenkins视图配置文件
COPY ./task/tasks.json /app/task/tasks.json
COPY --from=0 /app/devops-api ./
#动态文件，确保镜像变化
COPY --from=0 /app/.build_info ./
COPY --from=0 /app/config.docker.yaml ./

EXPOSE 5888
ENTRYPOINT /app/devops-api -c config.docker.yaml
