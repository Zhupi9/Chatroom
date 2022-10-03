package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)

func main() {
	//初始化redis链接池与userDao
	initPool("localhost:6379", 16, 0, time.Second*300)
	model.MyUserDao = model.NewUserDao(pool)
	//1."tcp"代表使用tcp协议
	//2."localhost:8888"代表监听本地8888端口
	l, err := net.Listen("tcp", "localhost:8972")
	if err != nil {
		fmt.Println("Listeing port 8972 error")
		return
	}
	fmt.Println("Server starts listening....")
	defer l.Close()
	var cliChan = make(chan bool, 0)
	//检查客户端的链接情况
	go func() {
		for {
			switch {
			case <-cliChan:
				fmt.Println("A client is disconnected")

			}
		}
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Connection Failed....")
		}
		fmt.Println("Connnection success, remote ip:", conn.RemoteAddr().String())

		go Processor(conn, cliChan) //开启一个协程处理事务,对应一个client

	}
}
