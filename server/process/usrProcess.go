// TODO 添加在用户登录时，返回当前在线用户列表
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
	//当前链接的用户名
	UserName string
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
	} else { //? 用户登录成功，101
		fmt.Printf("user %v Login ....\n", logInf.UserName)
		logRes.Code = message.LogSucc //101
		//?用户登录成功，将用户放到到OnlineUsers当中
		this.UserName = logInf.UserName
		userMgr.AddOnlineUser(this)
		//?将OnlineUSer添加到回复中发回客户端
		for i := range userMgr.onlineUsers {
			logRes.UserList = append(logRes.UserList, i)
		}

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
	//?如果登录成功101，则通知其他用户该用户已上线
	if logRes.Code == message.LogSucc {
		this.NotifyOthersOnline(this.UserName)
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
	} else { //?注册成功101
		fmt.Printf("user %v Login ....\n", RegInf.User.UserName)
		RegRes.Code = message.RegSucc //101
		//用户注册成功后直接登录，将用户放到到OnlineUsers当中
		this.UserName = RegInf.User.UserName
		userMgr.AddOnlineUser(this)
		//?将在线用户写到回复中
		for i := range userMgr.onlineUsers {
			RegRes.UserList = append(RegRes.UserList, i)
		}
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
	//?如果注册成功102，则通知其他用户该用户已上线
	if RegRes.Code == message.RegSucc {
		this.NotifyOthersOnline(this.UserName)
	}
	return
}

func (this *UsrProcess) NotifyOthersOnline(onlineName string) (err error) {
	//?定义发送消息，以及修改用户状态消息
	var mes message.Message
	var StaMes message.UserStatusMes

	StaMes.UserName = onlineName
	StaMes.Status = message.Online

	data, err := json.Marshal(StaMes)
	if err != nil {
		return
	}
	mes.Type = message.UserStatusMesType
	mes.Data = string(data)
	//?循环在线用户
	for name, link := range userMgr.onlineUsers {
		if name == this.UserName {
			continue
		}
		//?通知每个的用户
		link.NotifyMeOnline(mes)
	}
	return
}

// ?通知我，某用户上线了
func (this *UsrProcess) NotifyMeOnline(mes message.Message) (err error) {
	tf := &utils.Transfer{Conn: this.Conn}
	//?发送消息给客户端
	err = tf.WritePkg(mes)
	return
}
