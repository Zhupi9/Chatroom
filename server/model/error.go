package model

import "errors"

//根据业务逻辑需要，自定义一些错误

var (
	//?用户已存在（注册时）
	ERROR_USER_ALREADYEXIST = errors.New("用户已经存在")
	//?信息错误
	ERROR_USER_INFORMATION = errors.New("用户信息错误")
	//?其他错误（服务器端错误）
	ERROR_OTHER_SERVER = errors.New("服务器错误")
	//?用户不在线
	ERROR_USER_NOTONLINE = errors.New("当前用户不在线")
)
