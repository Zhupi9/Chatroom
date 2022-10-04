package main

import (
	"chatroom/client/process"
	"fmt"
	"net"
)

func main() {

	var choice string
	var isLoop bool = true
	conn, err := net.Dial("tcp", "localhost:8972") //连接后端服务器
	if err != nil {
		fmt.Println("ERROR 405: Can't link to server......")
		return
	}
	defer conn.Close()
	fmt.Println("Connected to Chatroom server!", conn)
	//新建一个process(client)实例
	var process = &process.UsrProcess{
		Conn:     conn,
		UserList: make(map[string]int, 1024),
	}

	for {
		fmt.Println("--------------欢迎登陆ChatRoom--------------")
		fmt.Println("\t\t\t1.登录系统")
		fmt.Println("\t\t\t2.注册账户")
		fmt.Println("\t\t\t3.退出系统")
		fmt.Println("输入（1～3）选择菜单：")
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			fmt.Println("Log in chatroom....")
			process.Login()
		case "2":
			fmt.Println("Sign in for a new user....")
			process.Register()
		case "3":
			fmt.Println("Exit Chatroom...")
			isLoop = false
		default:
			fmt.Println("Unexpected input!(1,2,3)")
		}

		if !isLoop {
			break
		}
	}
}
