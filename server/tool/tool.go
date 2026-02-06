package tool

import (
	"bufio"
	"chat/db"
	"chat/msg"
	"errors"
	"io"
	"log"
	"net"
)

// HandleClientMessage 处理客户端
func HandleClientMessage(conn net.Conn, hub *msg.Hub) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("server handleClientMessage panic recovered: %v\n", r)
		}
	}()
	reader := bufio.NewReader(conn)
	// 登录获取用户名
	username := handleRegisterOrLogin(reader, conn, hub)
	if username == "" {
		return
	}
	// 进入聊天室
	chatting(username, reader, conn, hub)
}

// handleRegisterOrLogin 处理登录注册的消息
func handleRegisterOrLogin(reader *bufio.Reader, conn net.Conn, room *msg.Hub) (username string) {
	for {
		initMsg, err := msg.ReadJsonMessage(reader)
		if err != nil {
			log.Printf("%s 在登录注册时失败", conn.RemoteAddr().String())
			return ""
		}
		initMsg.Conn = conn
		switch initMsg.Type {
		case msg.MessageRegister:
			// 注册只是把用户信息加入缓存
			msg.Register(initMsg)
			continue
		case msg.MessageJoin:
			status := room.Join(initMsg)
			if status {
				username = initMsg.Sender
				return username
			}
		default:
		}
	}

}

// chatting 处理登录注册之后的信息
func chatting(username string, reader *bufio.Reader, conn net.Conn, room *msg.Hub) {
	for {
		message, err := msg.ReadJsonMessage(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				//log.Printf("正常退出")
			} else {
				log.Printf("%s 被kill或者异常断开...\n", username)
				room.Leave(username)
			}
			return
		}
		message.Conn = conn
		switch message.Type {
		// 指令直接执行
		case msg.MessageLeave, msg.MessageList, msg.MessageRank, msg.MessageHeart:
			room.MsgChan <- message
		// 聊天信息存入stream
		default:
			// 聊天消息才异步入 Redis Streams
			_, err = db.AddStreamsData(message.Sender, message.Content, message.Receiver)
			if err != nil {
				log.Println("写入 Redis Streams 失败:", err)
			}
		}
	}
}
