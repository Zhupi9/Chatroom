package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf  [8192]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//var buf []byte = make([]byte, 1024)
	//读取发送长度
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		fmt.Println("read pks error")
		return
	}
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("read pkg body error")
		return
	}

	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.unmarshal error")
		return
	}
	return
}

func (this *Transfer) WritePkg(mes message.Message) (err error) {
	//先将message包装成[]byte
	data, err := json.Marshal(mes)
	if err != nil {
		return err
	}
	//发送长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)

	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		return
	}
	//发送data本身
	//fmt.Printf("data: %v\n", string(data))
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("this.Conn.write failed..")
		err = errors.New("Login Unsucessful because wrong name or pwd.")
		return
	}
	return

}
