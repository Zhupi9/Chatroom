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
	Conn net.Conn
	User user.User
}

// ?通知我，某用户上线了
func (this *UsrProcess) NotifyMeOnline(mes message.Message) (err error) {
	tf := &utils.Transfer{Conn: this.Conn}
	//?发送消息给客户端
	err = tf.WritePkg(mes)
	return
}

func (this *SmsProcess) SendSMSToOthers(mes message.Message) (err error) {

	var sms, smsRes message.SmsMes
	var res message.Message
	var data []byte = make([]byte, 1024)

	err = json.Unmarshal([]byte(mes.Data), &sms)
	if err != nil {
		return
	}
	smsRes.Content = sms.Content
	smsRes.SrcUser = sms.SrcUser
	res.Type = message.SmsMesType

	fmt.Println(sms.SrcUser, "send message to ", sms.DestUser)
	for _, v := range sms.DestUser {
		up, ok := userMgr.onlineUsers[v]
		if !ok {
			continue
		}
		//服务器端转发消息时目标客户端只有一个
		smsRes.DestUser = append(smsRes.DestUser, v)
		data, err = json.Marshal(smsRes)
		res.Data = string(data)
		err = up.GotSMS(res)
	}
	return
}

func (this *UsrProcess) GotSMS(mes message.Message) (err error) {
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(mes)

	return
}
