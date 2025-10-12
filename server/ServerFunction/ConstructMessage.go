package ServerFunction

import (
	"errors"
	"strings"
	"time"
)

// 这个包用来解析用户的命令
// 包括 register/login/private/list/exit

// ParseUserCmd 解析消息
func ParseUserCmd(input string) (*Message, error) {

	// 按照多个空格进行分割
	parts := strings.Split(input, " ")
	instructions := strings.ToUpper(parts[0])
	// 这里是否需要将 后面的消息进行拼接,因为消息内容也可以有空格
	arguments := parts[1:]
	switch instructions {
	//"LOGIN", "REGISTER"  这里是对用户消息的分析,为了防止用户在聊天过程中再次调用 登录/注册 这里不在能够识别登录注册
	case "PRIVATE":
		// 这种只能有两个参数 (账号,密码) (接收用户,私聊信息)
		if len(arguments) != 2 {
			return nil, errors.New("格式错误")
		}
	case "LIST", "QUIT":
		// 这种情况 不能有参数
		if len(arguments) != 0 {
			return nil, errors.New("不能有参数")
		}

	default:
		//普通广播消息 Broadcast
	}

	return &Message{
		Instructions: instructions,
		Arguments:    arguments,
		Text:         input,
		TimeStamp:    time.Now().Format("15:04:05"),
	}, nil

}

//// 如果是登录注册的内容
//cmd, username, password string, err error) {
//arr := strings.Split(msg, " ")
//if len(arr) != 3 {
//return "", "", "", errors.New("格式错误,应为: LOGIN username password或者 Register username password")
//}
//cmd = arr[0]
//username = arr[1]
//password = arr[2]
//if cmd != "LOGIN" && cmd != "REGISTER" {
//return "", "", "", errors.New("未知指令,只能使用 LOGIN或者 REGISTER")
//}
//return cmd, username, password, nil
