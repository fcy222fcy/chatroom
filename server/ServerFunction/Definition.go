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

// 所有的用户存放在一个map中---暂时不考虑在线问题

var OnlineUser = make(map[string]*Client)

var Mutex sync.Mutex

// 这里定义完管道之后,可以用于 分配路由的地方
var (
	// OnlineChan 加入聊天室提示---标识在线状态
	OnlineChan = make(chan string, 100)
	// 添加一个离开的

	// BroadcastChan 群聊
	BroadcastChan = make(chan *Message)
	// PrivateChan 私聊
	PrivateChan = make(chan *Message)
	// ListChan 列出名单
	ListChan = make(chan *Message)
	// HelpChan 帮助提示
	HelpChan = make(chan *Message)
)

// Message 结构体用来存放用户 消息
type Message struct {
	Instructions string   //分析出来的指令名称
	Arguments    []string //参数
	Text         string   //原始消息
	Sender       *Client  //发送者信息

	TimeStamp string
	IsOnline  bool //用户是否在线
}
