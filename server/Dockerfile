# 构建阶段
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./server

# 运行阶段
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
# 声明容器运行时监听的端口 8080
CMD ["./server"]
