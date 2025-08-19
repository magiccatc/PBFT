package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

// 客户端监听计时器，如果客户端在100ms内没有接收到反馈，则说明协议没有达成，数据要重新发送
func clientTcpListenTimer(timer *time.Timer) {
	for {
		fmt.Println("tcp.go clientTcpListenTimer是否还在工作")
		<-timer.C
		timer.Reset(time.Millisecond * resendDataTime)
		ChanNum <- "failget"
		fmt.Println("wo没有发送数据，计时器开始工作了")
	}
}

// 客户端使用的tcp监听,就是看你是否在客户端输入了信息
func clientTcpListen() {
	// 建立tcp服务
	listen, err := net.Listen("tcp", clientAddr)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()
	//var Mapmark map[string]int
	//timer := time.NewTimer(resendDataTime * time.Millisecond)
	//go clientTcpListenTimer(timer)
	for { // 无线循环，一直监听
		conn, err := listen.Accept() // 监听，等待客户端建立连接
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn) // 如果监听到内容，就读取出来，内容包括id，时间戳，和客户端地址等
		if err != nil {
			log.Panic(err)
		}
		//timer.Reset(time.Millisecond * resendDataTime)
		strB := string(b)
		fmt.Println("tcp.go  clientTcpListen" + strB)
		ChanNum <- strB
		if bStopMark {
			return
		}
	}

}

// 节点使用的tcp监听
func (p *pbft) tcpListen() {
	// func (p *pbft) tcpListen() {
	listen, err := net.Listen("tcp", p.node.addr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("节点开启监听，地址：%s\n", p.node.addr)
	defer listen.Close()

	i := 0
	for {
		conn, err := listen.Accept()

		i++
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}

		// writeData(db, strconv.Itoa(i), b)
		// fmt.Println("数据库顺序位i= ", i, "	，db的键key= ", strconv.Itoa(i))

		p.handleRequest(b)
	}

}

// 使用tcp发送消息
func tcpDial(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect error", err)
		return
	}

	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
