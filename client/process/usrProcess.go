package process

import (
	"chatroom/common/message"
	_ "chatroom/common/user"
	"chatroom/common/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UsrProcess struct {
	Conn     net.Conn
	UserName string
	UserList map[string]int
}

var (
	choice    string
	name, pwd string
	phone     string
)

// 客户端登录处理
func (this *UsrProcess) Login() (err error) {
	choice = "1"
	_, err = this.Conn.Write([]byte(choice))
	if err != nil {
		fmt.Println("Login Error...")
		return
	}

	//	准备发送和接收的消息
	var mes, res message.Message
	var logInf message.LoginMes
	var logRes message.LoginResMes
	var tf = &utils.Transfer{
		Conn: this.Conn,
	}

	mes.Type = message.LoginMesType

	var name, pwd string
	fmt.Print("Username:")
	fmt.Scanln(&name)
	fmt.Print("Password:")
	fmt.Scanln(&pwd)

	logInf = message.LoginMes{UserName: name, UserPwd: pwd}
	//将LoginRes Message序列化
	data, err := json.Marshal(logInf)
	if err != nil {
		return err
	}
	mes.Data = string(data)
	mes.Type = message.LoginMesType
	//发送登录信息给服务器端
	err = tf.WritePkg(mes)
	if err != nil {
		return
	}
	//读取服务器端的回复，code为202则登录成功，为404则登录失败
	res, err = tf.ReadPkg()

	json.Unmarshal([]byte(res.Data), &logRes)

	if logRes.Code == message.LogSucc {
		this.UserName = name
		//?显示当前在线用户,并添加到用户列表
		for i, v := range logRes.UserList {
			fmt.Println(i, ": ", v)
			this.UserList[v] = message.Online
		}
		//?登录成功，启动一个协程，保持与服务器的通讯
		//?如果服务器有消息推送给客户端，则接收他
		go this.serverProcessMes()
		this.ShowMenu()
	} else {
		fmt.Println(logRes.Error, logRes.Code)
	}

	return
}

func (this *UsrProcess) Register() (err error) {
	//告诉服务器要进行注册
	choice = "2"
	_, err = this.Conn.Write([]byte(choice))
	if err != nil {
		return
	}
	//准备发送和接收消息
	var mes, res message.Message
	var RegInf message.RegisterMes
	var RegRes message.RegisterResMes
	var tf = &utils.Transfer{
		Conn: this.Conn,
	}

	fmt.Print("Username:")
	fmt.Scanln(&name)
	fmt.Print("Password:")
	fmt.Scanln(&pwd)
	fmt.Print("PhoneNumber:")
	fmt.Scanln(&phone)

	RegInf.User.UserName = name
	RegInf.User.UserPwd = pwd
	RegInf.User.PhoneNumber = phone
	//将RegInf序列化
	data, err := json.Marshal(RegInf)
	if err != nil {
		return err
	}
	mes.Data = string(data)
	mes.Type = message.RegisterMesType
	//向服务器端发送注册请求消息
	err = tf.WritePkg(mes)
	if err != nil {
		return
	}

	//读取服务器端的回复，若为102则注册成功，403则用户名已存在，注册失败，505则服务器错误返回
	res, err = tf.ReadPkg()

	json.Unmarshal([]byte(res.Data), &RegRes)

	if RegRes.Code == message.RegSucc {
		this.UserName = name
		//?显示当前在线用户,并添加到用户列表
		for i, v := range RegRes.UserList {
			fmt.Println(i, ": ", v)
			this.UserList[v] = message.Online
		}
		//?注册成功，启动一个协程，保持与服务器的通讯
		go this.serverProcessMes()
		this.ShowMenu()
	} else {
		fmt.Println(RegRes.Error, RegRes.Code)
	}

	return
}
