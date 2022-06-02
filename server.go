package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播的channel
	Message chan string
}

// 创建一个server的接口
// create a server object
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		// 在线用户列表
		OnlineMap: make(map[string]*User),
		// 广播消息的channel
		Message: make(chan string),
	}
	return server
}

// 监听message广播消息channel的goroutine 一旦有消息发送给全部的用户
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		// 将消息发送给全部在线用户
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) handler(conn net.Conn) {
	// 处理当前连接的业务
	// fmt.Println("连接建立成功。。。")
	user := NewUser(conn, this)
	// 用户上线 加入OnlineMap
	user.Online()

	// 广播用户上线消息
	// this.BroadCast(user, "已上线")
	// 监听用户是否活跃的channel
	isLive := make(chan bool)
	// 接收客户端消息
	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if n == 0 {
				// this.BroadCast(user,"下线")
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				log.Println("conn read err", err)
				return
			}
			//将得到消息去除'\n'
			msg := string(buf[:n-1])
			// 得到消息进行广播
			// this.BroadCast(user, msg)
			//用户针对消息进行处理
			user.DoMessage(msg)

			// 表示当前用户依然活跃
			isLive <- true
		}

	}()
	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户是活跃的  重置定时器
			//不做任何事情 为了激活select 更新下面的定时器
		case <-time.After(10 * time.Second):
			//已经超时
			//将user强制下线
			user.SendMsg("你被踢了。。。")
			//销毁用的资源
			close(user.C)

			//关闭链接
			conn.Close()
			// user.Offline()
		
			//退出当前handler
			return // runtime.Goexit()
		}

	}
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
	// 启动监听message的goroutine
	go this.ListenMessager()
	for {
		// accept
		conn, err := listner.Accept()
		if err != nil {
			log.Println("accept err ...", err)
		}
		go this.handler(conn)
	}
}
