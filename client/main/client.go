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
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer ClientFunction.CloseConn(conn)

	// 连接上来之后直接让用户登录注册
	fmt.Println("欢迎来到Chat系统")
	fmt.Println("请先登录注册: LOGIN 用户名 密码 | REGISTER 用户名 密码")
	inputReader := bufio.NewReader(os.Stdin)
	// 这里需要改

	ClientFunction.LOGIN(conn)

	for {
		// 这里是登录完之后的菜单---这里还可以用原来的
		fmt.Println(" ============ 主界面 ============ ")
		fmt.Println("1.进入公共聊天室")
		fmt.Println("2.在公共聊天室输入LIST即可查看所有用户")
		fmt.Println("3.在公共聊天室输入[To]username message的格式即可私聊给对方")
		fmt.Println("4.退出系统---")
		fmt.Print("-请输入您的操作- ")
		input, errR := inputReader.ReadString('\n')
		if errR != nil {
			fmt.Println("errR :", errR)
			continue
		}
		input = strings.TrimSpace(input)
		switch input {
		//2和3还是基于1实现的
		case "2":
			fallthrough
		case "3":
			fallthrough
		case "1":
			//向服务端发送标记,防止用户注册完直接进入聊天室
			if err := Tool.Send(conn, []byte("_ENTER_")); err != nil {
				fmt.Println("发送标记err:", err)
				return
			}
			go ClientFunction.Receive(conn)
			fmt.Println("可以开始聊天了,输入exit/EXIT退出聊天室")
			for {
				if ClientFunction.Write(conn, inputReader) {
					break
				}
			}
		case "4":
			fmt.Println("正在退出")
			return
		default:
			fmt.Println("无效输入")
		}
	}
}
