package process

import (
	"chatroom/common/utils"
	"fmt"
	"net"
	"os"
)

func ShowMenu() {
	for {
		fmt.Println("----------------欢迎来到Chatroom---------------")
		fmt.Println("----------------1.在线用户--------------")
		fmt.Println("----------------2.发送消息--------------")
		fmt.Println("----------------3.消息记录--------------")
		fmt.Println("----------------4.退出系统--------------")
		fmt.Println("请选择（1-4）：")
		var choice string
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			fmt.Println("显示在线用户列表～")
		case "2":
			fmt.Println("发送消息～")
		case "3":
			fmt.Println("信息列表～")
		case "4":
			fmt.Println("你选择退出了系统～")
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新输入")
		}

	}

}

func serverProcessMes(conn net.Conn) {
	//创建一个transfer实例，不停的读取服务器发送的消息
	var tf = &utils.Transfer{
		Conn: conn,
	}
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			return
		}
		//根据读取到消息进行后面的处理
		fmt.Printf("mes: %v\n", mes)
	}
}
