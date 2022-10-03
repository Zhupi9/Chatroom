package process

import (
	"chatroom/common/message"
	"chatroom/common/user"
	"chatroom/common/utils"
	"chatroom/server/model"
	"encoding/json"
	"fmt"
	"net"
)

type UsrProcess struct {
	Conn net.Conn
}

func (this *UsrProcess) Login() (user *user.User, err error) {

	//服务器接受登录信息
	var mes, res message.Message
	var logInf message.LoginMes
	var logRes message.LoginResMes
	var tf = &utils.Transfer{
		Conn: this.Conn,
	}

	//读取客户端发来的登录信息
	mes, err = tf.ReadPkg()
	if err != nil {
		return
	}
	/*
		if mes.Type != message.LoginMesType { //如果读取到的消息类型不是登陆消息，则不受理
			fmt.Println("Not required type of message!")
			return
		}
	*/
	err = json.Unmarshal([]byte(mes.Data), &logInf)
	if err != nil {
		return
	}

	//判断登录用户名和密码，准备登录回复
	user, err = model.MyUserDao.Login(logInf.UserName, logInf.UserPwd)
	if err != nil {
		logRes.Error = err.Error()
		if err == model.ERROR_USER_INFORMATION {
			logRes.Code = message.LogFail //404
		} else {
			logRes.Code = message.ServerFail //505
		}
	} else {
		fmt.Printf("user %v Login ....\n", logInf.UserName)
		logRes.Code = message.LogSucc //101
	}
	fmt.Println(logRes.Error, logRes.Code)

	//将回复（LoginRes）序列化发送
	res.Type = message.LoginResMesType
	data, err := json.Marshal(logRes)
	if err != nil {
		fmt.Println("json.marshal failed...")
		return
	}
	res.Data = string(data)
	err = tf.WritePkg(res)
	if err != nil {
		return
	}

	return
}

func (this *UsrProcess) Register() (user *user.User, err error) {
	//服务器接受登录信息
	var mes, res message.Message
	var RegInf message.RegisterMes
	var RegRes message.RegisterResMes
	var tf = &utils.Transfer{
		Conn: this.Conn,
	}
	//接受客户端发来的注册信息
	mes, err = tf.ReadPkg()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(mes.Data), &RegInf)
	if err != nil {
		fmt.Println("json.unmarshal error")
		return
	}
	//将信息持久化到Redis中
	user, err = model.MyUserDao.Register(RegInf.User.UserName, RegInf.User.UserPwd, RegInf.User.PhoneNumber)
	//根据返回值确认是否持久化成功，准备返回值
	if err != nil {
		RegRes.Error = err.Error()
		if err == model.ERROR_USER_ALREADYEXIST {
			RegRes.Code = message.DupUser //403
		} else {
			RegRes.Code = message.ServerFail //505
		}
	} else {
		fmt.Printf("user %v Login ....\n", RegInf.User.UserName)
		RegRes.Code = message.RegSucc //101
	}
	fmt.Println(RegRes.Error, RegRes.Code)
	//将回复（RegRes）序列化发送
	res.Type = message.RegisterResMesType
	data, err := json.Marshal(RegRes)
	if err != nil {
		fmt.Println("json.marshal failed...")
		return
	}
	res.Data = string(data)
	err = tf.WritePkg(res)
	if err != nil {
		return
	}

	return
}
