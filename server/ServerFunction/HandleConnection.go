package ServerFunction

import (
	"log"
	"net"
)

// 处理连接,创建对应的结构体

func HandleConnection(conn net.Conn) *Client {
	client := &Client{
		Conn:     conn,
		Addr:     conn.RemoteAddr().String(),
		UserName: "",
		Status:   "默认在线",
	}
	log.Printf("新客户端连接: %s", client.Addr)

	// 不能够在这里关闭
	//defer func() {
	//	// 断开连接时清理用户
	//	if client.UserName != "" {
	//		Mutex.Lock()
	//		delete(OnlineUser, client.UserName)
	//		Mutex.Unlock()
	//	}
	//	err := conn.Close()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	log.Printf("客户端断开:%s", client.Addr)
	//}()

	return client
}
