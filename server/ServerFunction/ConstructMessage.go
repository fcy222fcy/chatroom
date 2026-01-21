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
	// parts := strings.Split(input, " ")
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, errors.New("空消息")
	}
	instructions := strings.ToUpper(parts[0])
	// 这里是否需要将 后面的消息进行拼接,因为消息内容也可以有空格
	arguments := parts[1:]

	switch instructions {
	//"LOGIN", "REGISTER"  这里是对用户消息的分析,为了防止用户在聊天过程中再次调用 登录/注册 这里不在能够识别登录注册
	case "PRIVATE":
		// 这种只能有两个参数 (账号,密码) (接收用户,私聊信息)
		if len(arguments) < 2 {
			return nil, errors.New("私聊格式错误:PRIVATE username message")
		}
		// 合并私聊消息内容--把合并的部分放在分析的地方,就是解析层和处理层之间选一个作为唯一责任方,不能两边都做
		username := arguments[0]
		msgText := strings.Join(arguments[1:], "")
		arguments = []string{username, msgText}
	case "LIST", "QUIT", "HELP":
		// 这种情况 不能有参数
		if len(arguments) != 0 {
			return nil, errors.New("命令不需要参数")
		}

	// 上面已经包含了这个参数
	case "_ENTER_", "_QUIT_":
		if len(arguments) != 0 {
			return nil, errors.New("系统命令不需要参数")
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
