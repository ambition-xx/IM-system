package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
		flag:       999,
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

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) PublicChat() {
	var charMsg string
	fmt.Println("请输入公聊消息 exit退出：")
	fmt.Scanln(&charMsg)

	for charMsg != "exit" {
		if len(charMsg) > 0 {
			_, err := client.conn.Write([]byte(charMsg + "\n"))
			if err != nil {
				fmt.Println("conn.write error :", err)
				break
			}
		}
		charMsg = ""
		fmt.Println("请输入公聊消息：")
		fmt.Scanln(&charMsg)
	}
}
func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scan(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.write error: ", err)
		return false
	}
	return true
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
			client.PublicChat()
		case 2:
			fmt.Println("私聊模式")
		case 3:
			client.UpdateName()
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
	go client.DealResponse()
	//  启动客户端业务
	client.Run()
	// select {}
}
