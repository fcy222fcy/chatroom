package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"time"
)

var DB *sqlx.DB

type User struct {
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

// GetUserByUsername 查询用户,包括 mysql和redis
func GetUserByUsername(username string) (*User, error) {
	// 旁路缓存
	// 读先读 redis,再读mysql
	// 写先写 mysql,再更新redis
	key := "user:" + username

	// 先查redis缓存
	userMap, err := RDB.HGetAll(key).Result()
	if err == nil && len(userMap) > 0 {
		// redis中有数据,手动赋值,不支持scan
		id, _ := strconv.Atoi(userMap["id"])
		user := &User{
			Id:       id,
			Username: userMap["username"],
			Password: userMap["password"],
		}
		if user.Username != "" {
			return user, nil
		}
	}

	// 2. 缓存未命中,查mysql
	var user User
	sqlStr := "select id,username,password from users where username = ?"
	err = DB.Get(&user, sqlStr, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("mysql search failed:%w", err)
	}
	// 写回 redis
	_, err = RDB.HMSet(key, map[string]interface{}{
		"id":       strconv.Itoa(user.Id),
		"username": user.Username,
		"password": user.Password,
	}).Result()
	if err != nil {
		// 缓存失败不影响主流程,记录日志
		log.Printf("failed to cache user %s:%v", username, err)
	} else {
		RDB.Expire(key, 30*time.Minute)
	}
	return &user, nil
}
