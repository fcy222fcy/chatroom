package ServerFunction

import (
	"chat/Tool"
	"fmt"
	"log"
	"strings"
	"time"
)

func (Client *Client) Process() {
	// 为每一个客户创建一个 context来控制生命周期
	//ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		// 断开连接时清理用户
		if Client.UserName != "" {

			delete(OnlineUser, Client.UserName)

		}
		err := Client.Conn.Close()
		if err != nil {
			log.Println(err)
		}
		log.Printf("客户端断开:%s", Client.Addr)
	}()

	// 简化Process()
	// 阶段1:用户身份认证
	if !Client.AuthenticationPhase() {
		return
	}

	// 测试--已运行
	log.Printf("用户已经登录成功,接下来等待用户进入聊天室信号")

	// 阶段2:等待用户进入聊天室信号
	if !Client.WaitForEnterSignal() {
		return
	}

	// 测试`-- 已执行
	//log.Printf("用户已经进入聊天室成功,接下来等待用户发送消息")

	// 阶段3:处理用户消息
	Client.MessageProcessingPhase()
}

// AuthenticationPhase 验证用户身份
func (Client *Client) AuthenticationPhase() bool {
	if err := Client.VerifyUser(); err != nil {
		log.Println("用户验证失败:", err)
		return false
	}
	return true
}

// WaitForEnterSignal 等待进入聊天室指令
func (Client *Client) WaitForEnterSignal() bool {
	// 测试
	//log.Printf("用户进入聊天室函数 WaitForEnterSignal")

	data, err := Tool.Recv(Client.Conn)
	if err != nil {
		log.Printf("读取_ENTER_信号失败:%v", err)
		return false
	}
	text := strings.TrimSpace(data)
	if text == "_ENTER_" {
		Client.EntryRoom()
		// 这里输出和客户端重复了
		//Tool.Send(Client.Conn, "欢迎进入聊天室,开始聊天\n")
		return true
	}
	return false
}

// MessageProcessingPhase 阶段3:处理用户消息
func (Client *Client) MessageProcessingPhase() {

	// 测试
	//log.Printf("用户开始发送消息")

	for {
		text, err := Tool.Recv(Client.Conn)
		if err != nil {
			log.Printf("读取用户消息失败:%v", err)
			break
		}

		// 测试信息
		//log.Printf("判断用户发送消息是否成功: %s", text)

		// 解析消息
		msg, err := ParseUserCmd(text)
		if err != nil {
			fmt.Println("解析消息失败:", err)
			continue
		}
		// 设置发送者
		msg.Sender = Client
		// 通过命令处理器分发消息
		msg.HandleCommand()
	}
}

func (Client *Client) EntryRoom() {

	// 首先先添加进入map
	Mutex.Lock()
	defer Mutex.Unlock()
	_, exists := OnlineUser[Client.UserName]
	// 如果不存没有的情况
	if exists {
		log.Printf("用户%s 已在线", Client.UserName)
		return
	}
	OnlineUser[Client.UserName] = Client

	OnlineChan <- &Message{
		//Instructions:    ,
		Text:      fmt.Sprintf("%s 进入了聊天室\n", Client.UserName),
		Sender:    Client,
		TimeStamp: time.Now().Format("15:04:05"),
	}
}
