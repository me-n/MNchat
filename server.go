package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type User struct {
	Conn net.Conn
	Name string
	Addr string
}

var (
	UserMap = make(map[string]User)
)

func SHandleError(err error, when string) {
	if err != nil {
		fmt.Println(err, when)
		os.Exit(1)
	}
}

//将客户端的信息保存在sever.log文件中——MADE IN MN
func sWriteLog(msg string) {
	fileName := "./sever.log"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	SHandleError(err, "os.OpenFile")
	defer file.Close()
	logMsg := time.Now().Format("2006-01-02 15:04:05") + msg
	file.Write([]byte(logMsg))
}

func HandConn(u User) {
	//通知其他客户端当前客户端上线了——MADE IN MN
	for _, u1 := range UserMap {
		//跳过通知自己——MADE IN MN
		if u1 != u {
			_, err := u1.Conn.Write([]byte(u.Name + "上线了\n"))
			SHandleError(err, "HandConn Conn.Write")
		}
	}
	defer u.Conn.Close()
	buf := make([]byte, 1024)
	//循环接收客户端信息——MADE IN MN
	for {
		n, e := u.Conn.Read(buf)
		if e == io.EOF {
			delete(UserMap, u.Name)
			fmt.Println(u.Name + "下线了")
			for _, u1 := range UserMap {
				_, e := u1.Conn.Write([]byte(u.Name + "下线了\n"))
				SHandleError(e, "HandConn u1.Conn.Write")
			}
			break
		}
		SHandleError(e, "HandConn conn.Read")
		if n > 0 {
			msg := string(buf[:n])
			//如果客户端发出退出指令，则断开此客户端连接——MADE IN MN
			if msg == "退出" {
				delete(UserMap, u.Name)
				fmt.Println(u.Name + "下线了")
				for _, u1 := range UserMap {
					_, e := u1.Conn.Write([]byte(u.Name + "下线了\n"))
					SHandleError(e, "HandConn u1.Conn.Write")
				}
				break
			}
			//判断是否是私聊——MADE IN MN
			if string(buf[:1]) == "@" {
				msg = string(buf[1:])
				split := strings.Split(msg, " ")
				if len(split) > 1 {
					acceptUser := split[0]
					msg = split[1]
					for key, user0 := range UserMap {
						if acceptUser == key {
							_, e := user0.Conn.Write([]byte(u.Name + "@ " + msg + "\n"))
							SHandleError(e, "HandConn user0.Conn.Write")
						}
					}
					sWriteLog(u.Name + "@ " + msg + "\n")
				}
			} else {
				//通知其他客户端当前客户端发送的信息（群聊）——MADE IN MN
				for _, u1 := range UserMap {
					if u1 != u {
						_, e := u1.Conn.Write([]byte(u.Name + ":" + msg + "\n"))
						SHandleError(e, "HandConn u1.Conn.Write")
					}
				}
				sWriteLog(u.Name + ":" + msg + "\n")
			}
		}
	}
}

func acceptConn(listen net.Listener) {
	var newUser User
	//循环接收客户端连接——MADE IN MN
	for {
		conn, err := listen.Accept()
		SHandleError(err, "listen.Accept")
		//接收客户端关于昵称的信息——MADE IN MN
		bufName := make([]byte, 256)
		n, err1 := conn.Read(bufName)
		if err1 != io.EOF {
			SHandleError(err1, "conn.Read user")
		}
		if n > 0 {
			userName := string(bufName[:n])
			//判断昵称是否存在——MADE IN MN
			if UserMap[userName].Name == userName {
				_, err := conn.Write([]byte("对不起，此昵称已存在！客户端已断开！"))
				SHandleError(err, "main conn.Write")
				continue
			}
			//新建客户端信息，并将其压入UserMap中——MADE IN MN
			newUser = User{conn, userName, conn.RemoteAddr().String()}
			UserMap[userName] = newUser
			fmt.Println(newUser.Name + "上线了")
		} else {
			//输入错误直接断开连接——MADE IN MN
			_, err := conn.Write([]byte("输入错误,连接已断开.\n"))
			SHandleError(err, "conn.Write")
			err2 := conn.Close()
			SHandleError(err2, "conn.Close")
			continue
		}
		//开独立协程处理此客户端交互——MADE IN MN
		go HandConn(newUser)
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:6666")
	SHandleError(err, "net.Listen")
	defer listen.Close()
	acceptConn(listen)
}
