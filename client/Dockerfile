# 构建阶段

# 命名为 builder
FROM golang:1.24 AS builder
# 容器内工作目录为 /app
WORKDIR /app
# 宿主机根目录所有文件复制到容器的 /app
# 第一个 . 表示 构建上下文的所有内容,也就是项目下的所有文件
# 第二个 . 表示 容器的/app目录
COPY . .
# 设置go 模块代理为国内镜像源,加速依赖下载
Run go env -w GOPROXY=https://goproxy.cn,direct

# 这一步的作用就是：1、整理模块依赖 2、将app/下的client包下的所有go文件编译成可执行的二进制文件存放在 /app/client
# 整理 Go 模块依赖，添加缺失的模块，移除未使用的模块
# CGO_ENABLED=0: 禁用 CGO，生成纯静态二进制文件，不依赖系统库
# GOOS=linux: 指定目标操作系统为 Linux
# GOARCH=amd64: 指定目标架构为 x86-64
# go build -o /app/client ./client
# -o /app/client: 指定输出文件路径和名称
# ./client: 要编译的包路径（通常是包含 main.go 的目录）
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/client ./client

# 运行阶段

# 基于 Alpine Linux 镜像，这是一个极简的 Linux 发行版
# 相比使用完整的 Golang 镜像，大幅减小最终镜像体积
FROM alpine:latest
# 设置运行阶段的工作目录为 /app
WORKDIR /app
#--from=builder: 从之前的构建阶段（builder）复制文件
#/app/client: 源文件路径（构建阶段生成的二进制文件）
#. : 目标路径（当前工作目录 /app）
COPY --from=builder /app/client .
#容器启动时执行的命令，执行当前目录下的 client 二进制文件
CMD ["./client"]



