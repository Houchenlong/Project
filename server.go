package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip string
	Port int

	//在线用户列表
	onlineMap map[string]*User
	mapLock	sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建一个server函数
func NewServer(ip string, port int) *Server{
	server := &Server{
		Ip:   ip,
		Port: port,
		onlineMap:make(map[string]*User),
		Message:make(chan string),
	}
	return server
}

//监听Message广播消息channel的goroutine
func(this *Server) ListenMessager(){
	for{
		msg := <- this.Message

		this.mapLock.Lock()
		defer this.mapLock.Unlock()
		for _, cli := range this.onlineMap{
			cli.C <- msg
		}
	}
}

//广播方法
func(this *Server) BoardCast(user *User, msg string){
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}


func (this *Server)Handler(conn net.Conn){

	user := NewUser(conn, this)
	//用户上线，将用户加入onlinemap中
	user.Online()


	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接受客户端消息
	go func() {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if n == 0{
			user.Offline()
			return
		}
		if err != nil && err != io.EOF{
			fmt.Println("conn read error:", err)
		}

		msg := string(buf[:n-1])

		//用户针对msg进行处理
		user.DoMessage(msg)

		//用户任意消息
		isLive <- true
	}()


	for{
		select {
		case <- isLive:
			//当前用户活跃，重置定时器
		case <- time.After(time.Second*10):
			//已经超时，将当前user踢掉
			user.SendMsg("你被踢了")
			//销毁资源
			close(user.C)
			//关闭连接
			conn.Close()
			//退出handler
			return
		}
	}

}

//启动服务器的方法
func (this *Server) Start(){
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil{
		fmt.Println("net.listen error:", err)
	}
	//close listen socket
	defer listener.Close()

	//启动监听message的goroutine
	go this.ListenMessager()

	for{
		//accept
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("conn accept error:", err)
		}
		//do  handler
		go this.Handler(conn)
	}
}