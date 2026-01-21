package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type user struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

// InitDB 连接数据库
func InitDB() (err error) {
	dsn := "root:123456@tcp(localhost:3306)/mychat?charset=utf8mb4&parseTime=True&loc=Local"
	//连接数据库并尝试ping
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect to mysql failed:%w", err)
	}
	return nil
}

// AddUserDb 注册用户
func AddUserDb(username string, password string) (err error) {
	sqlStr := "insert into users(username,password) values (?,?)"
	_, err = DB.Exec(sqlStr, username, password)
	if err != nil {
		return fmt.Errorf("mysql adduser failed:%w", err)
	}
	return nil
}

// 只是改了名字
// GetUserFromMySQL 查询用户
func GetUserFromMySQL(username string) (pwd string, err error) {
	var u user
	sqlStr := "select id,username,password from users where username = ?"
	err = DB.Get(&u, sqlStr, username)
	if err != nil {
		return "", fmt.Errorf("mysql search failed:%w", err)
	}
	return u.Password, nil
}
