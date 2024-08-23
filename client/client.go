package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
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

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("程序超时退出")
	}
}
