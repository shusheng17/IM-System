package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

// 用户上线业务
func (u *User) Online() {
	//用户上线，将用户加入OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//广播上线消息
	u.server.BroadCast(u, "online")
}

// 用户下线业务
func (u *User) Offline() {
	//用户下线，将用户从OnlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	//广播下线消息
	u.server.BroadCast(u, "offline")
}

// 向当前User对应的客户端发送消息
func (u *User) sendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "online...\n"
			u.sendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式： rename|张三
		newName := strings.Split(msg, "|")[1]
		// newName := msg[7:]	//可以不使用split

		//判断name是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.sendMsg("The name has been used\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.sendMsg("User name is updated: " + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式： to|张三|消息内容

		//1.获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.sendMsg("Message format is incorrect, please use \"to|jack|hello\" to send\n")
		}

		//2.根据用户名，获取对方User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.sendMsg("Username is not exist\n")
			return
		}

		//3.获取消息内容，通过User对象将消息发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.sendMsg("No message, try again\n")
			return
		}
		remoteUser.sendMsg(u.Name + "send:" + content + "\n")
	} else {
		//将得到的消息广播
		u.server.BroadCast(u, msg)
	}
}

// 监听当前User channel的方法，一旦有消息就发送给客户端
func (u *User) ListenMessage() {
	// for {
	// 	msg := <-u.C
	// 	u.conn.Write([]byte(msg + "\n"))
	// }

	//当u.C通道关闭后，不在进行（用于解决强踢后cpu飙升）
	for msg := range u.C {
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			panic(err)
		}
	}
}
