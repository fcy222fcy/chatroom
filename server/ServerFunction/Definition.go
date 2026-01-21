package ServerFunction

import (
	"net"
	"sync"
)

// Client 结构体用来存放用户 信息
type Client struct {
	UserName string   //用户名
	Conn     net.Conn //用户链接
	Status   string   //用户在线状态
	Addr     string
}

// 可以创建一个公共结构体--不能是常量,要不然值不能改变

// 所有的用户存放在一个map中---暂时不考虑在线问题

var OnlineUser = make(map[string]*Client)

var Mutex sync.Mutex

// 使用管道来缓冲,写成一个管道,然后将消息类型放在里面

// MessageChan 消息管道处理用户所有的,命令消息
var MessageChan = make(chan *Message, 1000)

// OnlineChan 在线通知管道,用于处理用户上线/下线通知
var OnlineChan = make(chan *Message, 1000)

// Message 结构体用来存放用户 消息
type Message struct {
	Instructions string   //分析出来的指令名称
	Arguments    []string //参数
	Text         string   //原始消息
	Sender       *Client  //发送者信息

	TimeStamp string
	IsOnline  bool //用户是否在线
}
