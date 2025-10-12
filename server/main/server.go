package main

import (
	"chat/server/ServerFunction"

	"log"
	"net"
)

func main() {
	// 初始化数据库
	err := ServerFunction.InitDB()
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 打开端口,启动监听
	Listen, errL := net.Listen("tcp", "localhost:8080")
	if errL != nil {
		log.Println("监听失败 :", errL)
		return
	}
	log.Println("监听已启动，等待连接接入")

	//直接使用defer 不能够处理err
	defer func(Listen net.Listener) {
		err := Listen.Close()
		if err != nil {
			log.Println("Error closing listener")
		}
	}(Listen)

	// 打开 广播连接

	//循环等待客户端连接
	for {
		conn, err := Listen.Accept()
		if err != nil {
			log.Println("connect err:", err)
		}
		// 传进去一个conn------转化为一个client结构体
		Client := ServerFunction.HandleConnection(conn)

		Client.Process()

	}
}
