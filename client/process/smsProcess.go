package process

import (
	"chatroom/common/message"
	"chatroom/common/user"
	"chatroom/common/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
	User user.User //当前客户端登录用户
	Conn net.Conn
}

func (this *UsrProcess) GotUsrStatusChange(mes message.Message) (err error) {
	var statusMes message.UserStatusMes
	err = json.Unmarshal([]byte(mes.Data), &statusMes)
	if err != nil {
		return
	}

	this.UserList[statusMes.UserName] = statusMes.Status
	fmt.Printf("%v上线啦！\n", statusMes.UserName)
	return
}

func (this *SmsProcess) SendSMS(content string, dest []string) (err error) {
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	//?准备发送SMS的消息
	var mes message.Message
	var smsMes message.SmsMes

	smsMes.Content = content
	smsMes.SrcUser = this.User.UserName
	smsMes.DestUser = dest
	data, err := json.Marshal(smsMes)
	if err != nil {
		return
	}
	mes.Type = message.SmsMesType
	mes.Data = string(data)

	err = tf.WritePkg(mes)

	return
}

func (this *SmsProcess) GotSMS(mes message.Message) (err error) {
	var sms message.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &sms)
	if err != nil {
		return
	}
	//?在客户端显示服务器转发的信息（包括发送来源和内容）
	fmt.Println(sms.SrcUser, "发送给你：", sms.Content)

	return
}
