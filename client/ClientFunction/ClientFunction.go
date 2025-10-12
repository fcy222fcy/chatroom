package ClientFunction

import (
	"bufio"
	"chat/Tool"
	"fmt"
	"net"
	"os"
	"strings"
)

var running = true

//func Login(conn net.Conn) {
//	// 这里 服务端要对应 ----- 只有接受到 才将登录注册进程结束
//	inputReader := bufio.NewReader(os.Stdin)
//	// 这里因为含有指令 发给服务端判断
//	text, _ := inputReader.ReadString('\n')
//	text = strings.TrimSpace(text)
//	_, err := conn.Write([]byte(text))
//	if err != nil {
//		log.Println("发送登录注册信息错误", err)
//	}
//	// 读取结果
//	_, err = conn.Read(make([]byte, 4))
//	if err != nil {
//		fmt.Println("接收服务器返回消息错误", err)
//	}
//
//}

// LOGIN 设置用户名
func LOGIN(conn net.Conn) {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		//fmt.Println("请输入您想设置的id")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)
		//校验空内容 如果用户输入回车,经过TrimSpace会变成""
		if text == "" {
			fmt.Println("输入为空,请重新输入!")
			continue
		}

		// 发送 登录/注册信息 给服务器
		if err := Tool.Send(conn, []byte(text)); err != nil {
			fmt.Println("发送信息出错:", err)
			return
		}

		// 接收服务器反馈---------------------------这里需要改进
		data, err := Tool.Recv(conn)
		if err != nil {
			fmt.Println("接收服务器反馈出错:", err)
			return
		}

		// 校验返回内容
		resp := strings.TrimSpace(string(data))
		if resp == "注册成功" {
			fmt.Println("注册成功！")
			break
		} else if resp == "登录成功" {
			fmt.Println("登录成功!")
			break
		} else {
			fmt.Println("登录注册失败了:", resp)
		}
	}
}

// CloseConn 关闭连接
func CloseConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println("关闭连接出错:", err)
	}
}

// Write 向服务端写入
func Write(conn net.Conn, inputReader *bufio.Reader) bool {
	input, _ := inputReader.ReadString('\n')
	//去掉首尾空白,包含\r\n,如果是回车和空格就会变成""
	text := strings.TrimSpace(input)
	if text == "" {
		fmt.Println("输入不能为空,请重新输入!")
		return false
	}
	//用户退出聊天室
	if strings.ToUpper(text) == "EXIT" {
		running = false
		return true
	}

	err := Tool.Send(conn, []byte(text))
	if err != nil {
		fmt.Println("发送消息出错:", err)
		return true
	}
	return false
}

// Receive 接收服务端返回的消息并输出
func Receive(conn net.Conn) {
	for running {
		data, err := Tool.Recv(conn)
		if err != nil {
			fmt.Println("与服务器的连接已断开", err)
			return
		}
		fmt.Println(string(data))
	}
}
