package ServerFunction

import (
	"chat/Tool"
	"log"
)

// 只分发,不做业务逻辑

func (msg *Message) HandleCommand() {
	// 确保用户在聊天时已经登录()

	if msg.Instructions != "LOGIN" && msg.Instructions != "Register" {
		Mutex.Lock()
		found := false
		for _, user := range OnlineUser {
			if msg.Sender == user {
				found = true
				break
			}
		}
		Mutex.Unlock()
		if !found {
			err := Tool.Send(msg.Sender.Conn, "请先登录\n")
			if err != nil {
				log.Println("登录出错", err)
			}
			return
		}
	}

	MessageChan <- msg
}

func StartCommandHandler() {
	for {
		select {
		case t := <-OnlineChan:
			text := t.Text
			HandleOnline(text)
		case msg := <-MessageChan:
			// 根据消息类型进行分发处理
			switch msg.Instructions {
			case "PRIVATE":
				msg.HandlePrivate()
			case "LIST":
				msg.HandleList()
			case "HELP":
				msg.HandleHelp()
			case "_QUIT_":
				msg.HandleQuit()
			default:
				msg.HandleBroadcast()
			}
		}
	}
}
