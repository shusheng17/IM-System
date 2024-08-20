package main

import "fmt"
import "net"
import "sync"

type Server struct{
	Ip string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建应该server接口
func NewServer(ip string, port int) *Server{
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

// 监听message广播消息channel的goroutine，一旦有消息就发送给所有在线的user
func (s *Server) ListenMessager(){
	for{
		msg := <-s.Message

		//将msg发送给所有在线的user
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息
func (s *Server) BroadCast(user *User, msg string){
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn){
	//当前链接的业务
	// fmt.Println("链接建立成功")

	user := NewUser(conn)

	//用户上线，将用户加入OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//广播上线消息
	s.BroadCast(user, "已上线")

	// 当前handle阻塞
	select {}
}

//启动服务器的接口
func (s *Server) Start(){
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil{
		fmt.Println("net.Listen err:", err)
		return
	}

	//close listen socket
	defer listener.Close()

	// 启动监听Message的gouroutine
	go s.ListenMessager()

	for{
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listenr accept err:", err)
			continue
		}

		//do handler
		go s.Handler(conn)

	}

	
}