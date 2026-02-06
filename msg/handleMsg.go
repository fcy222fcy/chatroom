package msg

import (
	"chat/db"
	"fmt"
	"log"
)

// HandleStreams 处理streams流消息
func (hub *Hub) HandleStreams() {
	lastID := "0-0"
	for {
		messages, err := db.ReadStreamsData(1, lastID)
		if err != nil {
			log.Println("读取 streams 出错:", err)
			continue
		}
		if len(messages) == 0 {
			continue
		}
		for _, m := range messages {
			sender := m.Values["sender"].(string)
			receiver := m.Values["receiver"].(string)
			content := m.Values["content"].(string)
			// 系统广播分支
			if sender == "系统广播" {
				hub.broadcast(receiver, fmt.Sprintf("%s: %s", sender, content))
				lastID = m.ID
				continue
			}

			msg := &Message{
				Sender:   sender,
				Receiver: receiver,
				Content:  content,
				Type:     MessageChat,
			}
			// 如果 sender 在线，再附加 Conn
			if client, ok := hub.Clients[m.Values["sender"].(string)]; ok {
				msg.Conn = client.Conn
			}
			if msg.Receiver != "" {
				hub.PrivateChat(msg)
			} else {
				hub.broadcast(msg.Sender, fmt.Sprintf("%s: %s", msg.Sender, msg.Content))
			}
			_ = db.AddActivity(msg.Sender, 1)
			lastID = m.ID // 更新游标，防止重复读取
		}
	}
}

// HandleChanMessages 普通消息处理
func (hub *Hub) HandleChanMessages() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server room.HandleMessages panic recovered: %v\n", err)
		}
	}()
	for msg := range hub.MsgChan {
		switch msg.Type {
		case MessageHeart:
			hub.PongHeart(msg.Sender)
		case MessageList:
			hub.List(msg.Sender, msg.Conn)
		case MessageLeave:
			hub.Leave(msg.Sender)
		case MessageRank:
			SendRank(msg.Sender, msg.Conn)
		default:
		}
	}
}
