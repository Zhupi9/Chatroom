package message

import "chatroom/common/user"

const (
	//?消息类型
	LoginMesType       = "LoginMes"
	LoginResMesType    = "LoginResMes"
	RegisterMesType    = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	UserStatusMesType  = "UserStatusMes"

	//?服务器响应码
	LogSucc    = 101 //登录成功
	LogFail    = 404 //登录失败，用户名或密码错误
	RegSucc    = 102 //注册成功
	DupUser    = 403 //用户已存在，注册失败
	ServerFail = 505 //服务器错误

	//?用户状态
	Online  = 1
	OffLine = 2
	Sleep   = 3
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMes struct {
	//UserId int `json:"userId"`
	UserName string `json:"username"`
	UserPwd  string `json:"userpwd"`
}

type LoginResMes struct {
	Code     int      `json:"code"`
	Error    string   `json:"error"`
	UserList []string `json:"userlist"`
}

type RegisterMes struct {
	User user.User `json:"user"`
}

type RegisterResMes struct {
	Code     int      `json:"code"`
	Error    string   `json:"error"`
	UserList []string `json:"userlist"`
}

// ?服务器端推送用户上线（离线）通知
type UserStatusMes struct {
	UserName string `json:"username"`
	Status   int    `json:"status"`
}
