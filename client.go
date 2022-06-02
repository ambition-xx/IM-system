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
	flag       int // 表示当前模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag: 999,
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

func (client *Client) Menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag < 4 { 
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字。。。")
		return false
	}
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式")
		case 2:
			fmt.Println("私聊模式")
		case 3:
			fmt.Println("更新用户名")
		case 0:
			fmt.Println("退出")

		}
	}
}
var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
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

	//  启动客户端业务
	client.Run()
	// select {}
}
