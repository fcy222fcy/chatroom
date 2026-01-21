package ClientFunction

import (
	"bufio"
	"chat/Tool"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var running = true

// OutoutChan 管道用于存放要输出的消息
var OutputChan = make(chan string, 100)

func StartOutputHandler() {
	go func() {
		for output := range OutputChan {
			fmt.Print(output)
		}
	}()
}

// LOGIN 设置用户名
func LOGIN(conn net.Conn) {

	inputReader := bufio.NewReader(os.Stdin)
	for {
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)

		//校验空内容 如果用户输入回车,经过TrimSpace会变成""
		if text == "" {
			fmt.Println("输入为空,请重新输入!")
			continue
		}

		if err := Tool.Send(conn, text+"\n"); err != nil {
			log.Println("发送登录注册信息错误", err)
			return
		}

		// 这里测试有没有读取到用户发出的内容
		//fmt.Println(text + "<-这是你输入的内容")

		// 接收服务器反馈
		data, err := Tool.Recv(conn)
		if err != nil {
			fmt.Println("接收服务器反馈出错:", err)
			return
		}
		resp := strings.TrimSpace(data)
		//校验返回内容
		if resp == "注册成功" {
			fmt.Println("注册成功！")
			break
		} else if resp == "登录成功" {
			fmt.Println("登录成功!")
			break
		} else {
			fmt.Println("失败了:", resp)
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
	if strings.ToUpper(text) == "_QUIT_" {
		err := Tool.Send(conn, "_QUIT_")
		if err != nil {
			fmt.Println("发送退出指令出错:", err)
		}
		running = false
		return true
	}

	err := Tool.Send(conn, text)
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
			OutputChan <- fmt.Sprintf("与服务器的连接已断开:%v\n", err)
			return
		}
		OutputChan <- data
	}
}
