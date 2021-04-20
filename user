package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name 	string
	conn 	net.Conn
	flag 	int
}

func NewClient(serverIp string,serverPort int) *Client{
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag: 		999,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil{
		fmt.Println("net.Dial error:", err)
	}
	client.conn = conn
	//返回对象
	return client
}

func (client *Client) UpdateName() bool{

	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil{
		fmt.Println("conn.Write error:", err)
		return false
	}

	return true
}


//处理server回应的消息
func (client *Client) DealResponse(){
	io.Copy(os.Stdout, client.conn)
}


func (client *Client) menu() bool{
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.修改用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3{
		client.flag = flag
		return true
	}else{
		fmt.Println(">>>>>请输入合法的数字")
		return false
	}
}

func (client *Client) Run(){
	for client.flag != 0{
		for client.menu() != true{
		}

		switch client.flag {
		case 1:
			fmt.Println("公聊模式选择...")
			client.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式选择...")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("修改用户名选择...")
			client.UpdateName()
			break
		}
	}
}

var (
	serverIp string
	serverPort int
)


func init(){
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口地址(默认8888)")
}

func main(){
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil{
		fmt.Println(">>>>>链接服务器失败...")
		return
	}


	go client.DealResponse()

	fmt.Println(">>>>>链接服务器成功...")

	client.Run()
}




func (client *Client) PublicChat(){
	var chatMsg string

	fmt.Println("请输入内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit"{
		if len(chatMsg) != 0{
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil{
				fmt.Println("conn write error:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers(){
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil{
		fmt.Println("conn write error:", err)
		return
	}
}

func (client *Client) PrivateChat(){
	var remoteName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名]， exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit"{
		fmt.Println(">>>>请输入，exit退出：")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write error:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>请输入，exit退出：")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象[用户名]， exit退出")
		fmt.Scanln(&remoteName)
	}
}
