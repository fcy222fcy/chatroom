package ServerFunction

import (
	"chat/Tool"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func (Client *Client) VerifyUser() error {

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()

	conn := Client.Conn

	//text, err := Tool.Recv(conn)
	//if err != nil {
	//	return err
	//}
	// 这里的代码把第一次输入的内容吞了
	//fmt.Println("收到用户验证消息:", text)
	//Tool.Send(conn, text+"\n")

	for {

		text, err := Tool.Recv(conn)
		if err != nil {
			// 返回结果待写
			return fmt.Errorf("读取用户消息失败: %v", err)
		}

		log.Println("收到用户验证消息:", text)
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "LOGIN") {
			username, password := Disassemble(text)
			// 检查该用户是否在线
			Mutex.Lock()
			if _, exists := OnlineUser[username]; exists {
				Mutex.Unlock()
				err := Tool.Send(conn, "用户已登录\n")
				if err != nil {
					log.Println(err)
				}
				continue

			}
			Mutex.Unlock()
			err := Login(username, password)
			if err != nil {
				err := Tool.Send(conn, "登录失败:"+err.Error()+"\n")
				if err != nil {
					log.Println(err)
				}
				continue
			}
			Client.UserName = username
			err = Tool.Send(conn, "登录成功\n")
			if err != nil {
				log.Println(err)
			}
			return nil
		} else if strings.HasPrefix(text, "REGISTER") {
			username, password := Disassemble(text)
			err := Register(username, password)
			if err != nil {
				err = Tool.Send(conn, "注册失败:"+err.Error()+"\n")
				if err != nil {
					log.Println(err)
				}
				continue
			}
			Client.UserName = username
			err = Tool.Send(conn, "注册成功\n")
			if err != nil {
				log.Println(err)
			}
			return nil
		} else {
			err := Tool.Send(conn, "您输入的格式不对,应为 LOGIN 用户名 密码 || REGISTER 用户名 密码")
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// Disassemble 用于分解参数
func Disassemble(str string) (res1, res2 string) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return "参数不足", "我也不知道"
	}
	return parts[1], parts[2]
}

// Register 注册
func Register(username, password string) error {
	//检查是否存在
	var tmp string
	err := DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&tmp)
	if err == nil {
		return errors.New("用户已经存在")
	}
	//if err != sql.ErrNoRows {
	//	return err
	//}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	// 插入新用户
	_, err = DB.Exec("INSERT INTO users(username,password) VALUES(?,?)", username, password)
	if err != nil {
		return err
	}

	// 这里不能这样写,不是用户登录完添加到 里面,应该是 进入聊天室才加入,或者是登录之后才添加

	// 将用户添加进入map----用不用再检验一遍username??
	//Mutex.Lock()
	//OnlineUser[username] = msg.Sender
	//Mutex.Unlock()
	return nil
}

// Login 登录
func Login(username, password string) error {
	var dbPwd string
	err := DB.QueryRow("SELECT password FROM users where username = ?", username).Scan(&dbPwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("用户不存在")
		}
		return err
	}

	if dbPwd != password {
		return errors.New("密码错误")
	}
	return nil
}
