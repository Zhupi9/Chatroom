package main

import (
	_ "chatroom/common/message"
	"chatroom/server/process"

	"fmt"
	"net"
)

// 循环接受客户端输入并进行处理
func Processor(conn net.Conn, bchan chan bool) {

	addr := conn.RemoteAddr().String()
	//新建一个process(server)实例

	defer func() {
		conn.Close()
		fmt.Println(addr, "Connection Closed")
	}()

	for {
		var buf []byte = make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			bchan <- true
			return
		}

		content := string(buf[:n])

		switch content {
		//处理登陆请求
		case "1":
			var process = &process.UsrProcess{
				Conn: conn,
			}
			fmt.Printf("Client %v want to log in\n", addr)
			_, err = process.Login()
			if err != nil {
				fmt.Println(err)
			}
		//处理注册请求
		case "2":
			var process = &process.UsrProcess{
				Conn: conn,
			}
			fmt.Printf("Client %v want to sign in for new user\n", addr)
			_, err = process.Register()
			if err != nil {
				fmt.Println(err)
			}
		case "3":
			fmt.Printf("Client %v exit\n", addr)
			return
		default:
		}
	}

}

/*
func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录消息

	case message.RegisterMesType:
		//处理注册消息
	}

	return nil
}
*/
