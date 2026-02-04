package ServerFunction

import (
	"log"
	"net"
)

// 处理连接,创建对应的结构体

func Connection(conn net.Conn) *Client {

	// 设置读写超时
	//conn.SetReadDeadline(time.Now().Add(30*time.Minute))
	//conn.SetWriteDeadline(time.Now().Add(10*time.Minute))

	client := &Client{
		Conn:     conn,
		Addr:     conn.RemoteAddr().String(),
		UserName: "",
		Status:   "默认在线",
	}
	log.Printf("新客户端连接: %s", client.Addr)
	return client
}
