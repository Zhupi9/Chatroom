package main

//TODO 修改服务器Process结构，根据接受消息类型作出反应

import (
	"chatroom/common/message"
	"chatroom/common/utils"
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
		//接受客户端发来消息
		tf := utils.Transfer{
			Conn: conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			return
		}
		err = serverProcessMes(conn, &mes)
		if err != nil {
			return
		}
	}

}

func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录消息
		var up = &process.UsrProcess{
			Conn: conn,
		}
		_, err = up.Login(*mes)
		if err != nil {
			return
		}
	case message.RegisterMesType:
		//处理注册消息
		var up = &process.UsrProcess{
			Conn: conn,
		}
		_, err = up.Register(*mes)
		if err != nil {
			return
		}
	case message.SmsMesType:
		//处理短消息
		var sp = &process.SmsProcess{
			Conn: conn,
		}
		sp.SendSMSToOthers(*mes)
		return
	}
	return
}
