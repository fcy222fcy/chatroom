package ServerFunction

import (
	//"context"
	"bufio"
	"fmt"
	"log"
)

func (Client *Client) Process() {
	// 这些内容也可以放在主函数当中,但是需要导入太多东西,也可以把主包去了,只留下主函数

	// 为每一个客户创建一个 context来控制生命周期
	//ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		// 断开连接时清理用户
		if Client.UserName != "" {
			Mutex.Lock()
			delete(OnlineUser, Client.UserName)
			Mutex.Unlock()
		}
		err := Client.Conn.Close()
		if err != nil {
			log.Println(err)
		}
		log.Printf("客户端断开:%s", Client.Addr)
	}()

	// 首先开始登录/注册
	go func() {
		Client.VerifyUser()

		// if err != nil{
		//	cancel()
		//}
	}()

	// 这里应该写如果出错进行的处理--下面还需要有,如果用户成功进入聊天室之后
	// 如果跟结构体进行绑定--原本的用户结构体创建应该放在登录/注册里面了
	go func() {
		Client.EntryRoom()
		//if err != nil{
		//	cancel()
		//}
		go func() {
			// 进入消息循环
			//Client.
		}()
	}()

	go StartCommandHandler()

	// 然后进入聊天室--这里接受的消息才会用ConstructMessage包下的解析函数

	// 读 客户端 消息循环
	reader := bufio.NewReader(Client.Conn)
	for {
		text, _ := reader.ReadString('\n')
		msg, err := ParseUserCmd(text)
		if err != nil {
			fmt.Println(err)
		}
		msg.Sender = Client

		BroadcastChan <- msg
	}

}

func (Client *Client) EntryRoom() {
	// 首先先添加进入map
	Mutex.Lock()
	_, ok := OnlineUser[Client.UserName]
	// 如果不存没有的情况
	if !ok {
		OnlineUser[Client.UserName] = Client
		OnlineChan <- fmt.Sprintf("%s 加入聊天室", Client.UserName)
	} else {
		// 这里应该让用户重新输入

	}
	Mutex.Unlock()
}
