package ServerFunction

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

//用于初始化数据库

// DB 是 MySQL 数据库连接 全局数据库
var DB *sql.DB

// InitDB 初始化
func InitDB() error {
	dsn := "root:123456@tcp(127.0.0.1:3306)/mychat?charset=utf8mb4&parseTime=true"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		// 这里返回错误,让主函数 决定是否退出
		return err
	}

	//测试连接
	err = DB.Ping()
	if err != nil {
		DB.Close()
		return fmt.Errorf("数据库连接失败:%v", err)
	}

	// 设置连接池参数，避免连接耗尽
	// 最大打开连接数
	DB.SetMaxOpenConns(100)
	// 最大空闲连接数
	DB.SetMaxIdleConns(80)
	// 连接最大生命周期1小时
	DB.SetConnMaxLifetime(1 * time.Hour)
	// 空闲连接超过30分钟
	DB.SetConnMaxIdleTime(30 * time.Minute)

	log.Println("✅️Mysql 连接成功")
	return nil
}
