package ServerFunction

// 这里是所有的业务逻辑
import (
	"fmt"
	"log"
	"net"
)

func HandleOnline(text string) {

	// 这里其实还是用了广播的概念,主要是向别的用户发送...已上线的标识

	Mutex.Lock()
	for _, client := range OnlineUser {
		_, err := client.Conn.Write([]byte(text))
		if err != nil {
			log.Println(err)
		}
	}
	Mutex.Unlock()

}
func (msg *Message) HandleList() {
	conn := msg.Sender.Conn

	_, err := conn.Write([]byte("用户在线列表如下"))
	if err != nil {
		log.Println("发送用户在线列表失败:", err)
		return
	}

	Mutex.Lock()
	for name, client := range OnlineUser {
		userInfo := fmt.Sprintf("- %s (%s)\n", name, client.Addr)
		_, err = client.Conn.Write([]byte(userInfo))
		if err != nil {
			log.Println("发送用户信息失败", err)
			continue
		}
	}
	Mutex.Unlock()
	_, err = conn.Write([]byte("列表结束\n"))
	if err != nil {
		log.Println("发送列表结束标记失败:", err)
	}
}

func (msg *Message) HandlePrivate() {
	// 发送者的名字msg.Sender.UserName
	// 通过接受者的名字(msg.Arguments[1])找到对应的 结构体实例
	// 参数2是消息本体 msg.Arguments[2]

	// 添加时间戳
	timeStamp := msg.TimeStamp

	// 这里是否也可以用数据库处理
	// 查找用户conn
	// 创建一个接收者的 conn
	var conn net.Conn
	// 发送者的conn
	var FromConn net.Conn = msg.Sender.Conn
	// 创建一个标志,标志是否查找到
	var flag bool = false
	// 使用map,这里加锁
	Mutex.Lock()
	for username := range OnlineUser {
		if username == msg.Arguments[1] {
			conn = OnlineUser[username].Conn
			flag = true
			break
		}
	}
	Mutex.Unlock()

	// 处理标志
	if flag {
		// 给目标用户
		text := fmt.Sprintf("[%s] %s 给您发来一条消息:%s\n", timeStamp, msg.Sender.UserName, msg.Arguments[2])
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Println(err)
		}
	} else {
		// 给原来用户
		text := fmt.Sprintf("[%s] : 未找到该用户", timeStamp)
		_, err := FromConn.Write([]byte(text))
		if err != nil {
			log.Println(err)
		}
	}
}

// HandleQuit 是从聊天室中退出的,不是从菜单界面退出的
func (msg *Message) HandleQuit() {
	// 从map中将用户删除,就是退出了 EnterRoom
	timestamp := msg.TimeStamp
	username := msg.Sender.UserName
	Mutex.Lock()
	for _, client := range OnlineUser {
		if username == client.UserName {
			// 将用户从map中删除
			delete(OnlineUser, username)
		}
	}
	Mutex.Unlock()
	text := fmt.Sprintf("[%s][系统]:%s退出了聊天室", timestamp, username)
	// 系统日志
	log.Printf("[%s] 用户%s 加入聊天室成功\n", timestamp, username)
	// 对全用户通知
	SystemMessageBroadcast(text)

}

// HandleEnterRoom 这里才将用户添加进入聊天室
func (msg *Message) HandleEnterRoom() {
	timestamp := msg.TimeStamp

	username := msg.Sender.UserName
	client := msg.Sender
	// 这里才将用户添加进入聊天室
	Mutex.Lock()
	OnlineUser[username] = client
	Mutex.Unlock()
	text := fmt.Sprintf("[%s][系统]:%s加入了聊天室", timestamp, username)
	// 服务端输出
	log.Printf("[%s] 用户%s 退出聊天室\n", timestamp, username)
	// 客户端通知
	SystemMessageBroadcast(text)
}

// HandleBroadcast 处理用户消息的广播
func (msg *Message) HandleBroadcast() {
	timestamp := msg.TimeStamp
	username := msg.Sender.UserName
	// 取原始消息
	text := msg.Text
	// 包装消息
	text = fmt.Sprintf("[%s] %s : %s", timestamp, username, text)
	Mutex.Lock()
	for _, client := range OnlineUser {
		conn := client.Conn
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Println("广播发送错误", err)
		}
	}
	Mutex.Unlock()
}

// SystemMessageBroadcast 用于将系统消息对用户进行广播
func SystemMessageBroadcast(text string) {

	Mutex.Lock()
	for _, client := range OnlineUser {
		conn := client.Conn
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Println(text)
		}
	}
	Mutex.Unlock()
}

func (msg *Message) HandleHelp() {
	conn := msg.Sender.Conn
	text := fmt.Sprintf("这里是帮助说明:\n" +
		"1.登录注册:LOGIN USERNAME PASSEORD\n" +
		"2.查看在线用户列表:LIST\n" +
		"3.私聊其他用户PRIVATE USERNAME TEXT\n" +
		"4.进入聊天室_ENTER_\n" +
		"5.退出聊天室_QUIT_\n")
	_, err := conn.Write([]byte(text))
	if err != nil {
		log.Println("HELP 出错:", err)
	}
}
