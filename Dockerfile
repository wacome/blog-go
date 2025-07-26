FROM golang:1.23.0-alpine AS builder

# 设置国内 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go mod 和 go sum 文件
COPY go.mod go.sum ./

# 下载依赖项
# 这一步会利用 Docker 的缓存机制，只有在 go.mod 或 go.sum 变更时才会重新下载
RUN go mod download

# 复制所有源代码到工作目录
COPY . .

# 构建 Go 应用
# -ldflags="-w -s" 用于减小二进制文件大小
# CGO_ENABLED=0 确保静态链接，以便在 Alpine 镜像中运行
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/blog-go main.go


# ---- Final Stage ----
# 使用更小的 Alpine 镜像作为最终的运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/blog-go /app/blog-go

# 创建用于文件上传的目录
RUN mkdir -p /app/uploads/avatars \
             /app/uploads/images \
             /app/uploads/books \
             /app/uploads/friend-avatars \
             /app/uploads/comment-avatars

# 暴露应用端口（在 compose.yml 中会映射到主机端口）
EXPOSE 2233

# 容器启动时运行的命令
CMD ["/app/blog-go"]