package chatroom

import (
	"fmt"
	"net"
)

func Room(conn net.Conn) {
	var inbuf string
	fmt.Println("-------------------------聊天室--------------------------")
	for {
		fmt.Println("输入：")
		fmt.Scanln(&inbuf)
		if inbuf == "exit" {
			return
		}
		fmt.Println(inbuf)
	}
}
