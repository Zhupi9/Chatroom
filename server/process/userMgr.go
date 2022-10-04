package process

import "chatroom/server/model"

type UserMgr struct {
	onlineUsers map[string]*UsrProcess
}

//?因为usermgr在全服务器有且只有一个
//?并且在很多地方都用得到，因此定义为全局变量

var (
	userMgr *UserMgr
)

func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[string]*UsrProcess, 1024),
	}
}

func (this *UserMgr) AddOnlineUser(up *UsrProcess) {
	this.onlineUsers[up.UserName] = up
}

func (this *UserMgr) DelOnlineUser(up *UsrProcess) {
	delete(this.onlineUsers, up.UserName)
}

func (this *UserMgr) GetAllOnlineUSer() map[string]*UsrProcess {
	return this.onlineUsers
}

func (this *UserMgr) GetUPByName(name string) (up *UsrProcess, err error) {
	up = &UsrProcess{}
	up, ok := this.onlineUsers[name]
	if !ok {
		//? 说明当前用户不在线
		err = model.ERROR_USER_NOTONLINE
		return
	}
	return
}
