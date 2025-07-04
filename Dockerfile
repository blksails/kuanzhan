# 多阶段构建 Dockerfile 用于 kuanzhan CLI 工具

# 使用官方 golang 镜像作为构建阶段
FROM golang:1.23-alpine AS builder

# 安装必要的工具
RUN apk update && apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kuanzhan ./cmd/kuanzhan

# 使用轻量级镜像作为运行时阶段
FROM alpine:latest

# 安装 ca-certificates 以支持 HTTPS 请求
RUN apk --no-cache add ca-certificates

# 创建非 root 用户
RUN addgroup -g 1000 kuanzhan && \
    adduser -D -s /bin/sh -u 1000 -G kuanzhan kuanzhan

# 设置工作目录
WORKDIR /home/kuanzhan

# 从构建阶段复制二进制文件
COPY --from=builder /app/kuanzhan .

# 创建配置目录
RUN mkdir -p .config/kuanzhan

# 复制配置文件示例
COPY kuanzhan.yaml.example .config/kuanzhan/

# 更改文件所有者
RUN chown -R kuanzhan:kuanzhan /home/kuanzhan

# 切换到非 root 用户
USER kuanzhan

# 设置入口点
ENTRYPOINT ["./kuanzhan"]

# 默认命令
CMD ["--help"]

# 添加标签
LABEL org.opencontainers.image.title="Kuanzhan CLI"
LABEL org.opencontainers.image.description="快站管理命令行工具 - 快速创建、管理和部署快站站点"
LABEL org.opencontainers.image.vendor="Your Organization"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/your-org/kuanzhan" 