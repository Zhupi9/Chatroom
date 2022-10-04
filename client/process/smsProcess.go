package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

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
