FROM golang:alpine as builder

WORKDIR /app
# 复制源代码并编译
COPY . .
# 程序编译打包 设置环境变量
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://proxy.golang.org,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o devops-api .

FROM alpine:latest

LABEL MAINTAINER="David588@gmail.com"

WORKDIR /app

COPY --from=0 /app/devops-api ./
# COPY --from=0 /go/src/github.com/flipped-aurora/gin-vue-admin/server/resource ./resource/
COPY --from=0 /app/config.docker.yaml ./

EXPOSE 5888
ENTRYPOINT /app/devop-api -c config.docker.yaml
