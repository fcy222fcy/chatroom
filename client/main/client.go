package main

import (
	"bufio"
	"chat/Tool"
	"chat/client/ClientFunction"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 启动输出处理器
	go ClientFunction.StartOutputHandler()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer ClientFunction.CloseConn(conn)

	// 连接上来之后直接让用户登录注册
	ClientFunction.OutputChan <- "欢迎来到Chat系统\n"
	ClientFunction.OutputChan <- "请先登录注册: LOGIN 用户名 密码 | REGISTER 用户名 密码\n"
	ClientFunction.OutputChan <- ">"

	ClientFunction.LOGIN(conn)

	for {
		// 这里是登录完之后的菜单---这里还可以用原来的
		ClientFunction.OutputChan <- " ============ 主界面 ============ \n"
		ClientFunction.OutputChan <- "1.进入公共聊天室\n"
		ClientFunction.OutputChan <- "2.退出系统\n"
		ClientFunction.OutputChan <- ">"

		inputReader := bufio.NewReader(os.Stdin)
		input, errR := inputReader.ReadString('\n')

		if errR != nil {
			ClientFunction.OutputChan <- fmt.Sprintln("读取输入错误 :", errR)
			continue
		}
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			//向服务端发送标记,防止用户注册完直接进入聊天室
			if err := Tool.Send(conn, "_ENTER_"); err != nil {
				ClientFunction.OutputChan <- fmt.Sprintln("进入聊天室失败:", err)
				continue
			}

			// 启动接收协程--如果接收到数据就打印
			go ClientFunction.Receive(conn)
			// 这里能不能在用户首次进入聊天室的时候,发送一条帮助信息
			// 其他帮助信息......
			ClientFunction.OutputChan <- "****************************\n"
			ClientFunction.OutputChan <- "** 欢迎进入聊天室\n"
			ClientFunction.OutputChan <- "** 可用命令:\n"
			ClientFunction.OutputChan <- "**  LIST - 查看在线用户列表\n"
			ClientFunction.OutputChan <- "**  [PRIVATE]username message - 私聊指定用户\n"
			ClientFunction.OutputChan <- "**  HELP - 查看详细帮助信息\n"
			ClientFunction.OutputChan <- "**  exit/EXIT - 退出聊天室\n"
			ClientFunction.OutputChan <- "****************************\n"
			ClientFunction.OutputChan <- "可以开始聊天了,输入exit/EXIT退出聊天室\n"

			for {
				if ClientFunction.Write(conn, inputReader) {
					break
				}
			}
		case "2":
			ClientFunction.OutputChan <- "正在退出\n"
			return

		default:
			ClientFunction.OutputChan <- "无效输入\n"
		}
	}
}
