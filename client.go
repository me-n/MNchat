package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	quitChan=make(chan bool,0)
)

func CHandleError(err error, when string) {
	if err != nil {
		fmt.Println(err, when)
		os.Exit(1)
	}
}

func help() {
	help := "欢迎使用MN聊天室，单聊请输入'@'+'昵称'+'空格'+'内容'（例：@小明 你好)；查看聊天记录请输入'@log' ；查看帮助文档输入'@help'；结束聊天请输入'退出'。 "
	fmt.Println(help)
}
//向服务器发送信息——MADE IN MN
func cWrite(dial net.Conn,userName string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		lineByte, _, e := reader.ReadLine()
		if e != io.EOF {
			CHandleError(e, "reader.ReadLine")
		}else if e==io.EOF{
			continue
		}
		//设置帮助文档，发送"help"获取——MADE IN MN
		if string(lineByte) == "@help" {
			help()
		}else if string(lineByte) == "@log"{
			openLog(userName)
		}else{
			_, e = dial.Write(lineByte)
			CHandleError(e, "cWrite dial.Write")
			cWriteLog(string(lineByte),userName)
		}
	}
}
//接收服务器的信息——MADE IN MN
func cRead(dial net.Conn,userName string) {
	buf := make([]byte, 1024)
	for {
		n, err := dial.Read(buf)
		if err != io.EOF {
			CHandleError(err, "dial.Read")
		}
		if n > 0 {
			fmt.Print(string(buf[:n]))
			cWriteLog(string(buf[:n]),userName)
		}
	}
}
//将聊天记录写入本地以自己昵称命名的log文件中——MADE IN MN
func cWriteLog(msg string,u string)  {
	filePath:="./"+u+".log"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	CHandleError(err,"write os.OpenFile")
	defer file.Close()
	logMsg:=time.Now().Format("2006-01-02 15:04:05")+msg+"\n"
	file.Write([]byte(logMsg))
}
//读取本地聊天记录——MADE IN MN
func openLog(u string)  {
	filePath:="./"+u+".log"
	open, err := os.Open(filePath)
	CHandleError(err,"log os.Open")
	defer open.Close()
	buf:=make([]byte,2048)
	n, err1 := open.Read(buf)
	if err1!=io.EOF{
		CHandleError(err1,"log os.Read")
	}
	if n>0{
		fmt.Print(string(buf[:n]))
	}
}
/*
此版本为0。0001版，包含基础的群聊、私聊、聊天日志保存、聊天日志读取等基本功能，还十分初级，待完善。仅用于制作人应聘所用；
制作人：孟宁（林）
*/
func main() {
	//上线前设置昵称并上传服务器——MADE IN MN
	fmt.Println("请输入昵称：")
	reader := bufio.NewReader(os.Stdin)
	nameByte, _, e1 := reader.ReadLine()
	if e1 != io.EOF {
		CHandleError(e1, "read nameByte")
	}
	dial, err := net.DialTimeout("tcp", "127.0.0.1:6666",time.Second*5)
	CHandleError(err, "net.DialTimeout")
	_, e1 = dial.Write(nameByte)
	CHandleError(e1, "dial.Write")
	defer dial.Close()
	fmt.Println("欢迎使用MN聊天室，单聊请输入'@'+'昵称'+'空格'+'内容'（例：@小明 你好)；查看聊天记录请输入'@log' ；查看帮助文档输入'@help'；结束聊天请输入'退出'。 ")
	//循环写入信息——MADE IN MN
	go cWrite(dial,string(nameByte))
	//循环接收服务器传来的信息——MADE IN MN
	go cRead(dial,string(nameByte))
	<-quitChan
}
