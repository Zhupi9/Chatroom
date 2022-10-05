package process

import (
	"bufio"
	"chatroom/common/message"
	"chatroom/common/utils"
	"fmt"
	"os"
	"strings"
)

func (this *UsrProcess) ShowMenu() {
	sms := &SmsProcess{
		Conn: this.Conn,
		User: this.User,
	}
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
		case "1": //?显示所有用户（包括离线用户）
			this.ShowUserList()
		case "2":
			var destUsrName, content string
			var destUsrList []string = make([]string, 0)
			fmt.Println("选择发送消息对象们:(输入End以结束)")
			for {
				fmt.Scanln(&destUsrName)
				if strings.EqualFold(destUsrName, "end") {
					break
				}
				destUsrList = append(destUsrList, destUsrName)
			}
			fmt.Println("输入发送消息内容:")
			content, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				fmt.Println(err)
			} else {
				err = sms.SendSMS(content, destUsrList)
				if err != nil {
					fmt.Println(err)
				}
			}

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

func (this *UsrProcess) serverProcessMes() {
	//创建一个transfer实例，不停的读取服务器发送的消息
	var tf = &utils.Transfer{
		Conn: this.Conn,
	}
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			return
		}
		//根据读取到消息进行后面的处理
		switch mes.Type {
		case message.UserStatusMesType:
			err = this.GotUsrStatusChange(mes)
			if err != nil {
				fmt.Println(err)
			}
		case message.SmsMesType:
			sp := &SmsProcess{
				Conn: this.Conn,
			}
			err = sp.GotSMS(mes)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
