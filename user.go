package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	
	// 当前user属于那个server
	server *Server
}

// create a user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 监听当前user 里面的channel
	go user.ListenMessage()

	return user
}

//用户上线业务
func (this *User) Online() {
	// 用户上线 加入OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播用户上线消息
	this.server.BroadCast(this, "已上线")

}

//用户下线业务
func (this *User) Offline() {
	// 用户下线 从OnlineMap中删除user
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播用户下线消息
	this.server.BroadCast(this, "下线")
}

func (this *User)SendMsg(msg string){
	this.conn.Write([]byte(msg))
}
// 用户处理消息业务
func (this *User) DoMessage(msg string) {
	if msg == "who"{
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap{
			onilneMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线。。。\n"
			this.SendMsg(onilneMsg)

		} 
		this.server.mapLock.Unlock()
	} else {

		this.server.BroadCast(this, msg)
	}
	
}


func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}

