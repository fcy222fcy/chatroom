package msg

import (
	"chat/db"
	"chat/utils"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net"
	"time"
)

// broadcast 广播
func (hub *Hub) broadcast(sender, content string) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	temp := &Message{
		Type:    MessageChat,
		Sender:  sender,
		Content: content,
	}

	for username, client := range hub.Clients {
		if username == sender {
			continue
		}

		err := SendJsonMessage(client.Conn, temp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(sender, "已退出聊天室...连接已关闭")
				return
			}
			log.Println("broadcast:", err)
			return
		}
	}
	fmt.Println(content)
}

// PrivateChat 私聊
func (hub *Hub) PrivateChat(msg *Message) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	target, ok := hub.Clients[msg.Receiver]
	// 用户不在在线列表里面
	if !ok {
		_ = SendJsonMessage(msg.Conn, &Message{
			// 为什么这里写成了 广播??
			Type:    MessageChat,
			Sender:  "[系统]",
			Content: fmt.Sprintf("用户 %s 不存在或不在线", msg.Receiver),
		})
		return
	}
	err := SendJsonMessage(target.Conn, &Message{
		Type:    MessagePrivate,
		Sender:  msg.Sender,
		Content: msg.Content,
	})
	if err != nil {
		log.Println("PrivateChat:", err)
		return
	}
	fmt.Printf("%s 私聊 %s: %s\n", msg.Sender, msg.Receiver, msg.Content)
}

// List 查询在线列表
func (hub *Hub) List(name string, conn net.Conn) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	list := "在线用户列表: "
	for username := range hub.Clients {
		list += username + "  "
	}

	err := SendJsonMessage(conn, &Message{
		Type:    MessageList,
		Content: list,
	})
	if err != nil {
		log.Println("List ", err)
	}
	fmt.Println(name, "请求查看用户列表...")
}

// Register 处理注册信息
func Register(msg *Message) {
	err := db.AddUserDb(msg.Sender, msg.Content)
	if err != nil {
		// 检查是否是唯一约束冲突（用户名已存在）
		if isDuplicateKeyError(err) {
			rr := SendJsonMessage(msg.Conn, &Message{
				Type:    MessageRegister,
				Content: "用户名: " + msg.Sender + " 已被注册",
			})
			if rr != nil {
				log.Println("Register send error:", rr)
			}
		} else {
			log.Println("注册失败:", err)
			rr := SendJsonMessage(msg.Conn, &Message{
				Type:    MessageRegister,
				Content: "注册失败，请稍后重试",
			})
			if rr != nil {
				log.Println("Register send error:", rr)
			}
		}
		return
	}
	user := &db.User{
		Username: msg.Sender,
		Password: msg.Content,
	}
	setErr := db.SetUserToRedis(user)
	if setErr != nil {
		log.Println("将用户信息写入 Redis 失败:", setErr)
	}

	// 注册成功
	rr := SendJsonMessage(msg.Conn, &Message{
		Type:    MessageRegister,
		Content: "OK",
	})
	if rr != nil {
		log.Println("Register send error:", rr)
	}
	fmt.Println(msg.Sender, "注册成功...")
}

// isDuplicateKeyError 辅助函数检查是否是唯一约束错误
func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 // 1062是MySQL的重复键问题
	}
	return false
}

// Join 处理登录消息
func (hub *Hub) Join(msg *Message) bool {
	// 获取用户 redis+mysql
	user, err := db.GetUserByUsername(msg.Sender)
	// 查询失败的情况
	if err != nil {
		var respContent string
		if err.Error() == "user not found" {
			respContent = fmt.Sprintf("%s 不存在，请先注册", msg.Sender)
		} else {
			respContent = "登录失败，数据库异常"
			log.Printf("查询用户 %s 失败: %v", msg.Sender, err)
		}
		// 发送错误响应
		if r := SendJsonMessage(msg.Conn, &Message{
			Type:    MessageChat,
			Content: respContent,
		}); r != nil {
			log.Println("发送登录失败响应错误:", err)
		}
		return false
	}
	// 判断密码
	if user.Password != msg.Content {
		if err1 := SendJsonMessage(msg.Conn, &Message{
			Type:    MessageChat,
			Content: "密码错误，请重新输入",
		}); err1 != nil {
			log.Println("发送密码错误响应错误:", err)
		}
		return false
	}
	// 判断是否已经登录
	if _, ok := hub.Clients[msg.Sender]; ok {
		if r := SendJsonMessage(msg.Conn, &Message{
			Type:    MessageChat,
			Content: "该账户已登录",
		}); r != nil {
			log.Println("发送账号已登陆响应错误:", err)
		}
		return false
	}

	// 登录成功
	rr := SendJsonMessage(msg.Conn, &Message{
		Type:    MessageRegister,
		Content: "OK"})
	if rr != nil {
		log.Println("Register send error:", rr)
		return false
	}

	client := &Client{Username: msg.Sender, Conn: msg.Conn, LastActive: time.Now()}
	hub.AddClient(msg.Sender, client)
	//content := fmt.Sprintf("系统广播：%s 加入了聊天室...", msg.Sender)
	//hub.broadcast(msg.Sender, content)

	// 发送历史消息,可以返回一个切片或者map,取不同的字段
	historyMsg, rrr := db.ShowHistory()
	if rrr != nil {
		log.Println(rrr)
	}
	message := &Message{
		Type:    MessageChat,
		Content: historyMsg,
	}
	r := SendJsonMessage(msg.Conn, message)
	if r != nil {
		log.Println("发送历史消息失败:", r)
	}

	// 加入streams流
	_, err = db.AddStreamsData("系统广播", fmt.Sprintf("%s 加入了聊天室...", msg.Sender), msg.Sender)
	if err != nil {
		log.Println("写入 Redis Streams 失败:", err)
	}
	// 登录 增加2活跃度
	err = db.AddActivity(msg.Sender, 2)
	if err != nil {
		log.Println(msg.Sender, "登录增加活跃度失败 :", err)
	}
	return true
}

// Leave 处理退出消息
func (hub *Hub) Leave(username string) {
	_, err := db.AddStreamsData("系统广播", fmt.Sprintf("%s 离开了聊天室...", username), username)
	if err != nil {
		log.Println("Leave写入 Redis Streams 失败:", err)
	}
	hub.RemoveClient(username)
}

// PongHeart 处理心跳
func (hub *Hub) PongHeart(username string) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()
	if client, exists := hub.Clients[username]; exists {
		client.LastActive = time.Now()
		err := client.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			log.Printf("PongHeart: %v", err)
		}
	}
}

// StartHeartbeatMonitor 服务端定期检测客户端心跳超时
func (hub *Hub) StartHeartbeatMonitor() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		now := time.Now()
		hub.Mutex.Lock()
		for username, client := range hub.Clients {
			if now.Sub(client.LastActive) > 20*time.Second {
				log.Printf("用户 %s 心跳超时，强制下线\n", username)
				utils.CloseConn(client.Conn, username)
				hub.Leave(username)
			}
		}
		hub.Mutex.Unlock()
	}
}

// SendRank 发送活跃度排行
func SendRank(username string, conn net.Conn) {
	sprintf, err := db.ShowActivityRank()
	if err != nil {
		log.Println(err)
		return
	}
	rr := SendJsonMessage(conn, &Message{Type: MessageRank, Content: sprintf})
	if rr != nil {
		log.Printf("向%s发送活跃度排名失败:%s", username, rr)
		return
	}
	fmt.Println(username, "查看活跃度排行...")
}
