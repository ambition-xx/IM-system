package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

// create a server object
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) ListenMessager(){
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap{ 
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}
func (this *Server) handler(conn net.Conn) {
	// 处理当前连接的业务
	// fmt.Println("连接建立成功。。。")

	// 用户上线 加入OnlineMap
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播用户上线消息
	this.BroadCast(user, "已上线")
}

// start a server
func (this *Server) Start() {
	//socket lisen
	listner, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		log.Println("listen err ...", err)
	}
	// close socket
	defer listner.Close()
	go this.ListenMessager()
	for {
		// accept
		conn, err := listner.Accept()
		if err != nil {
			log.Println("accept err ..", err)
		}
		go this.handler(conn)
	}
}

