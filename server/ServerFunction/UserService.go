package ServerFunction

// 这里是所有的业务逻辑
import (
	"chat/Tool"
	"fmt"
	"log"
)

func HandleOnline(text string) {

	// 这里其实还是用了广播的概念,主要是向别的用户发送...已上线的标识

	Mutex.Lock()
	for _, client := range OnlineUser {
		err := Tool.Send(client.Conn, text)
		if err != nil {
			log.Println(err)
		}
	}
	Mutex.Unlock()

}
func (msg *Message) HandleList() {
	conn := msg.Sender.Conn
	err := Tool.Send(conn, "用户列表如下\n")
	if err != nil {
		log.Println("发送列表开始标记失败:", err)
		return
	}

	Mutex.Lock()
	defer Mutex.Unlock()
	for name, client := range OnlineUser {
		userInfo := fmt.Sprintf("- %s (%s)\n", name, client.Addr)
		err := Tool.Send(conn, userInfo)
		if err != nil {
			log.Println("循环发送用户在线列表失败", err)
			continue
		}
	}
	err = Tool.Send(conn, "列表结束\n")
	if err != nil {
		log.Println("发送列表结束标记失败:", err)
	}
}

func (msg *Message) HandlePrivate() {
	// 发送者的名字msg.Sender.UserName
	// 通过接受者的名字(msg.Arguments[1])找到对应的 结构体实例
	// 参数2是消息本体 msg.Arguments[2]

	// 参数边界检查
	if len(msg.Arguments) != 2 {
		err := Tool.Send(msg.Sender.Conn, "私聊格式错误\n")
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 添加时间戳
	timeStamp := msg.TimeStamp
	targetUsername := msg.Arguments[0]
	PrivateMessage := msg.Arguments[1]

	// 使用map,这里加锁
	Mutex.Lock()
	targetClient, exists := OnlineUser[targetUsername]
	Mutex.Unlock()

	if !exists {
		text := fmt.Sprintf("[%s] 系统:未找到该用户 %s\n", timeStamp, targetUsername)
		err := Tool.Send(msg.Sender.Conn, text)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 找到该用户
	targetText := fmt.Sprintf("[%s] %s给您发来一条私信: %s \n", timeStamp, msg.Sender.UserName, PrivateMessage)
	err := Tool.Send(targetClient.Conn, targetText)
	if err != nil {
		log.Println("发送私聊消息失败:", err)
	}

	// 给发送者确认
	//senderText := fmt.Sprintf("[%s] [私聊->%s]: %s\n", timestamp, targetUsername, privateMessage)
	//Tool.Send(msg.Sender.Conn, senderText)
}

// HandleQuit 是从聊天室中退出的,不是从菜单界面退出的
func (msg *Message) HandleQuit() {
	// 从map中将用户删除,就是退出了 EnterRoom
	timestamp := msg.TimeStamp
	username := msg.Sender.UserName

	Mutex.Lock()
	delete(OnlineUser, username)
	Mutex.Unlock()
	text := fmt.Sprintf("[%s][系统]:%s退出了聊天室\n", timestamp, username)
	// 系统日志
	log.Printf("[%s] 用户%s 加入聊天室成功\n", timestamp, username)
	// 对全用户通知
	SystemMessageBroadcast(text)
	err := Tool.Send(msg.Sender.Conn, "您已退出聊天室,返回主菜单\n")
	if err != nil {
		log.Println(err)
	}
}

// HandleBroadcast 处理用户消息的广播
func (msg *Message) HandleBroadcast() {
	timestamp := msg.TimeStamp
	username := msg.Sender.UserName
	// 取原始消息
	text := msg.Text
	// 包装消息
	text = fmt.Sprintf("[%s] %s : %s\n", timestamp, username, text)
	Mutex.Lock()
	defer Mutex.Unlock()
	for _, client := range OnlineUser {
		// 不向自己广播
		if client != msg.Sender {
			err := Tool.Send(client.Conn, text)
			if err != nil {
				log.Println("广播发送错误", err)
			}
		}
	}
}

// SystemMessageBroadcast 用于将系统消息对用户进行广播
func SystemMessageBroadcast(text string) {

	Mutex.Lock()
	for _, client := range OnlineUser {
		conn := client.Conn
		err := Tool.Send(conn, text)
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
		"5.退出聊天室_QUIT_\n" +
		"使用示例:\n" +
		"-发送公共消息:直接输入消息内容\n" +
		"-私聊:PRIVATE 张三 你还好吗?\n" +
		"=======================\n")
	err := Tool.Send(conn, text)
	if err != nil {
		log.Println("HELP 出错:", err)
	}
}
