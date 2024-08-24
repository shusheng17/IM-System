package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

// 处理server回应的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	//一旦client.conn有数据，姐直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改用户名")
	fmt.Println("4.显示在线用户")
	fmt.Println("0.退出")

	fmt.Scan(&flag)

	if flag >= 0 && flag <= 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法的数字>>>>>")
		return false
	}
}

func (client *Client) ShowUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.ShowUsers()
	fmt.Println(">>>>>请输入聊天对象[用户名],exit退出:")
	fmt.Scanln(&remoteName) //清除缓存"\n"
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>请输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>请输入聊天内容,exit退出")
			fmt.Scanln(&chatMsg)
		}

		client.ShowUsers()
		fmt.Println(">>>>>请输入聊天对象[用户名],exit退出:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>请输入聊天内容,exit退出")
	fmt.Scanln(&chatMsg) //清除缓存"\n"
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) Update() bool {
	fmt.Println(">>>>>请输入用户名：")
	fmt.Scan(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
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
			client.PublicChat()
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式")
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名")
			client.Update()
			break
		case 4:
			//显示在线用户
			fmt.Println("显示在线用户")
			client.ShowUsers()
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

	go client.DealResponse()

	client.Run()

	// select {
	// case <-time.After(10 * time.Second):
	// 	fmt.Println("程序超时退出")
	// }
}
