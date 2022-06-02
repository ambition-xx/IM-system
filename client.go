package main

import (
	"flag"
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

var serverIp string
var serverPort int
// ./client -ip 127.0.0.1 -port 8888
func init(){
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口号(默认是8888)")
}
func main() {

	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		log.Println("服务器链接失败。。。")
	}
	fmt.Println("服务器链接成功。。。")

	select {}
}
