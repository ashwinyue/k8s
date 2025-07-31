# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 整理模块依赖
RUN go mod tidy

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 创建日志目录
RUN mkdir -p /app/logs && chown -R appuser:appgroup /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV LOG_LEVEL=info
ENV APP_VERSION=1.0.0
ENV GIN_MODE=release

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./main"]