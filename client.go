package main

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		log.Println("net.dial error...", err)
	}
	client.conn = conn
	// 返回对象
	return client
}
func main() {
	client := NewClient()
	if client == nil {
		log.Println("服务器链接失败。。。")
	}
	fmt.Println("服务器链接成功。。。")

	select {}
}
