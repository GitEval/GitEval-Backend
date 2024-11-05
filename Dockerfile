# 第一阶段：构建 Go 应用
FROM golang:alpine AS builder

# 设置 Go 代理为七牛云的代理
ENV GOPROXY=https://goproxy.cn,direct

# 安装需要的依赖
RUN apk update && apk add --no-cache git

# 切换到 be-patent 目录，构建二进制文件
WORKDIR /app

# 复制 be-patent 服务代码
COPY . .

RUN go mod tidy && go build -o gitEval

# 第二阶段：复制编译结果到最终镜像
FROM alpine

# 设置工作目录为
WORKDIR /app

# 从 builder 复制编译好的二进制文件
COPY --from=builder /app/gitEval /app/gitEval


# 开放端口（根据需要设置）
EXPOSE 8080

# 启动用户服务
CMD ["./gitEval","-conf","conf/config.yaml"]
