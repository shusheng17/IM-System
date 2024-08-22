package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

//创建一个用户的API
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

//用户上线业务
func (u *User) Online() {
	//用户上线，将用户加入OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//广播上线消息
	u.server.BroadCast(u, "online")
}

//用户下线业务
func (u *User) Offline() {
	//用户下线，将用户从OnlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	//广播下线消息
	u.server.BroadCast(u, "offline")
}

//向当前User对应的客户端发送消息
func (u *User) sendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

//用户处理消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "online...\n"
			u.sendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else {
		//将得到的消息广播
		u.server.BroadCast(u, msg)
	}
}

//监听当前User channel的方法，一旦有消息就发送给客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
