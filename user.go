package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C	chan string
	conn net.Conn
	server *Server
}

//创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User{
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server:server,
	}

	//启动监听user channel消息的goroutine
	go user.ListenMessage()

	return user
}

//用户上线
func(this *User) Online(){
	//用户上线，将用户加入onlinemap中
	this.server.mapLock.Lock()
	defer this.server.mapLock.Unlock()
	this.server.onlineMap[this.Name] = this
	//广播用户上线
	this.server.BoardCast(this, "已上线")
}
//用户下线
func(this *User) Offline(){
	//用户下线，将用户加入onlinemap中
	this.server.mapLock.Lock()
	defer this.server.mapLock.Unlock()
	delete(this.server.onlineMap, this.Name)
	//广播用户下线
	this.server.BoardCast(this, "下线")
}

func (this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

//处理业务
func(this *User) DoMessage(msg string){
	if msg == "who"{
		//查询当前用户都有哪些
		this.server.mapLock.Lock()
		defer this.server.mapLock.Unlock()
		for _, user := range this.server.onlineMap{
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.SendMsg(onlineMsg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|"{
		//消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]

		//判断name是否被占用
		if _, ok := this.server.onlineMap[newName]; ok{
			this.SendMsg("用户名已存在")
			return
		}else{
			this.server.mapLock.Lock()
			delete(this.server.onlineMap, this.Name)
			this.server.onlineMap[newName] = this
			defer this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("您已经更新用户名"+this.Name+"\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|"{
		//消息格式： to|张三|消息内容

		//1. 获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == ""{
			this.SendMsg("消息格式不正确，请使用\"to|张三|你好\"格式。\n")
		}
		//2. 根据用户名 得到对方User对象
		remoteUser, ok := this.server.onlineMap[remoteName]
		if !ok{
			this.SendMsg("该用户不存在\n")
			return
		}
		//3. 根据消息内容，通过对方的User对象将消息内容发送
		content := strings.Split(msg, "|")[2]
		if content == ""{
			this.SendMsg("请重发")
		}
		remoteUser.SendMsg(this.Name+ "说："+ content)
	}
	this.server.BoardCast(this, msg)
}


//监听当前user channel的方法
func (this *User) ListenMessage(){
	for{
		msg := <- this.C
		this.conn.Write([]byte(msg+"\n"))
	}
}