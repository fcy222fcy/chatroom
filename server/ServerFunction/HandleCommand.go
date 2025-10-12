package ServerFunction

// 只分发,不做业务逻辑

func (msg *Message) HandleCommand() {
	// 首先检查用户是否在... 不然就返回错误提示
	Mutex.Lock()
	var flag bool
	for _, user := range OnlineUser {
		if msg.Sender == user {
			flag = true
			break
		}
	}
	if flag {
		//return "用户还未登陆--"
	}
	// 如果没找到进行下一步

	Mutex.Unlock()
	switch msg.Instructions {
	case "PRIVATE":
		PrivateChan <- msg
	case "LIST":
		ListChan <- msg
	case "HELP":
		HelpChan <- msg
	case "_ENTER_":
		msg.HandleEnterRoom() // 进入聊天室的逻辑直接处理
	case "_QUIT_":
		msg.HandleQuit() // 退出聊天室的逻辑直接处理
	default:
		BroadcastChan <- msg // 默认按广播处理
	}
}

func StartCommandHandler() {
	for {
		select {
		//这一条是向其他用户发送...已上线
		case text := <-OnlineChan:
			HandleOnline(text)
		// 专门用于广播/群聊
		case msg := <-BroadcastChan:
			msg.HandleBroadcast()
		case msg := <-PrivateChan:
			msg.HandlePrivate()
		case msg := <-ListChan:
			msg.HandleList()
		case msg := <-HelpChan:
			msg.HandleHelp()
		}
	}
}
