package model

import (
	"chatroom/common/user"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// 定义一个UserDao结构体
// 完成对User结构体的各种操作
var MyUserDao *UserDao

type UserDao struct {
	pool *redis.Pool
}

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	fmt.Println("UserDao created...")
	return
}

// 根据用户id返回User实例
func (this *UserDao) getUserByName(conn redis.Conn, name string) (usr *user.User, err error) {
	//在redis中查询用户名
	err = ERROR_OTHER_SERVER
	var res string
	usr = &user.User{}
	res, err = redis.String(conn.Do("HGET", "users", name))
	if err != nil {
		fmt.Println("getuser by name error...", err)
		if err == redis.ErrNil {
			err = ERROR_USER_INFORMATION
		}
		return
	}
	//将res反序列化成user对象
	err = json.Unmarshal([]byte(res), usr)
	if err != nil {
		fmt.Println("json.unmarshal failed ", err)
		return
	}

	err = nil //顺利执行，error置为空
	return
}

func (this *UserDao) saveUser(conn redis.Conn, name, pwd, phone string) (usr *user.User, err error) {
	err = ERROR_OTHER_SERVER
	//将信息序列化
	usr = &user.User{
		UserName:    name,
		UserPwd:     pwd,
		PhoneNumber: phone,
	}
	data, err := json.Marshal(usr)
	if err != nil {
		fmt.Println("json.marshal error", err)
		return
	}

	//将user存储到Redis中
	_, err = conn.Do("HSET", "users", name, data)
	if err != nil {
		return
	}
	err = nil
	return
}

// 供外部函数调用
func (this *UserDao) Login(name, pwd string) (usr *user.User, err error) {
	conn := this.pool.Get()

	defer conn.Close()
	usr, err = this.getUserByName(conn, name)
	if err != nil {
		return
	}

	//用户名存在，校验密码
	if usr.UserPwd != pwd {
		err = ERROR_USER_INFORMATION
		return
	}

	return
}

func (this *UserDao) Register(name, pwd, phone string) (usr *user.User, err error) {
	conn := this.pool.Get()

	defer conn.Close()
	//先查看Redis中是否已经存在该用户名
	usr, err = this.getUserByName(conn, name)
	if err == nil {
		err = ERROR_USER_ALREADYEXIST
		return
	}
	//用户名没有重复，可以注册
	usr, err = this.saveUser(conn, name, pwd, phone)
	if err != nil {
		return
	}

	return
}
