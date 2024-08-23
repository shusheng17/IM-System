package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
	}

	//连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	//返回对象
	return client
}

func (client *Client) menu() bool {
	var flag int
	// flag = -1

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改用户名")
	fmt.Println("0.退出")

	fmt.Scan(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法的数字>>>>>")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {

		}

		//根据不同模式处理不同业务
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式")
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式")
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名模式")
			break
		}
	}
}

var serverIp string
var serverPort int

// .\client.exe -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置默认服务器IP(默认为127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置默认服务器端口(默认为8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>连接失败>>>>>")
		return
	}
	fmt.Println(">>>>>连接成功>>>>>")

	client.Run()

	// select {
	// case <-time.After(10 * time.Second):
	// 	fmt.Println("程序超时退出")
	// }
}
