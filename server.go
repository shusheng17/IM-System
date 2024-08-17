package main

import "fmt"
import "net"

type Server struct{
	Ip string
	Port int
}

//创建应该server接口
func NewServer(ip string, port int) *Server{
	server := &Server{
		Ip: ip,
		Port: port,
	}

	return server
}

func (s *Server) Handler(conn net.Conn){
	//当前链接的业务
	fmt.Println("链接建立成功")
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