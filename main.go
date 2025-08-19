package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const nodeCount = 50 // 共识节点数，与main.go/pbftserver保持一致
const sleepTime = 5 // 一次共识结束，随机0-sleepTime ms的时间来发送下一组共识数据，
//随机休眠的最大时间为5ms，与main.go/pbftserver保持一致

// 客户端的监听地址
var clientAddr = "127.0.0.1:8888"
var ChanNum = make(chan string, nodeCount*100) //获取得到的确认的节点数
var bStopMark = false                          // 是否停止发送数据
var chanStop = make(chan bool, 1)              // 主程序是否终止
var stopNum = 3                                // 达成多少次共识后主程序退出，var stopNum = 10

// 节点池，主要用来存储监听地址
var nodeTable map[string]string // 其实就是
var startTime time.Time

const counttime = 300         // 计时，单位为秒，程序统计多长时间的数据
const resendDataTime = 100000 // 如果200ms没有收到数据，则任务共识没有达成，重新发送数据，const resendDataTime = 200000
const timeResentWait = 1000   // 这个必须小于上面那个数值，因为电脑算力有限，所以节点数越多，这个值就要越大，这里是在34个节点时我测试的不死的，最后得结果要减去这个值

// var endTime time.Time
// 根据节点数自动生成端口号和对应的IP地址，对应下面的nodeTable = map[string]string
func initNodeTable() {
	nodeTable = make(map[string]string, 2)
	for i := 0; i < nodeCount; i++ {
		numberS := fmt.Sprintf("N%d", i) //numberS=N0,N1,N2
		IPs := ""
		if i < 99 {
			IPs = fmt.Sprintf("127.0.0.1:80%02d", i) //IPs= 127.0.0.1:8000,1,2
		} else {
			IPs = fmt.Sprintf("127.0.0.1:8%02d", i) //IPs= 127.0.0.1:8000,1,2
		}

		nodeTable[numberS] = IPs
	}
}

// 解析节点信息，返回节点地址映射
func parseNodeInfo(nodeInfo string) map[string]string {
	nodes := make(map[string]string)
	for _, pair := range strings.Split(nodeInfo, ",") {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			nodes[parts[0]] = parts[1]
		}
	}
	return nodes
}

func main() {

	// // 创建一个通道来控制主程序的结束
	// chanStop := make(chan bool)
	// defer close(chanStop)

	//为四个节点生成公私钥，并把信息保存在keys目录下
	genRsaKeys()
	fmt.Println("公钥和私钥生成成功")

	// nodeTable = map[string]string{
	// 	"N0": "127.0.0.1:8000",
	// 	"N1": "127.0.0.1:8001",
	// 	"N2": "127.0.0.1:8002",
	// 	"N3": "127.0.0.1:8003",
	// 	"N4": "127.0.0.1:8004",
	// 	"N5": "127.0.0.1:8005",
	// 	"N6": "127.0.0.1:8006",
	// 	// "N120": "127.0.0.1:8120",
	// }
	initNodeTable()
	if len(os.Args) != 2 { //os.Args的类型是 []string ，也就是字符串切片。用来获取命令行参数可以用 len(os.Args) 来获取其数量.
		log.Panic("输入的参数有误！")
	}
	nodeID := os.Args[1]

	if nodeID == "client" {
		clientSendMessageAndListen() //启动客户端程序，当发送完客户请求信息，其实该线程的任务就完成了，这一步发出的消息是tcpDial(content, nodeTable["N0"])
		// 发送停止信号，通知主程序客户端已经发送完5次请求
		// chanStop <- true
		// <-chanStop
		// fmt.Println("客户端已发送完5次请求，主程序即将结束。")
	} else if addr, ok := nodeTable[nodeID]; ok { // nodeID就是自己输入的N0，N1……，addr就是对应的IP和端口号；addr, ok := nodeTable[nodeID]这句话中如果找到nodeID对应的信息，这ok为true，否则为false

		p := NewPBFT(nodeID, addr) // 创建PBFT信息，这是一个结构体，里面保存共识需要的所有信息
		// db, err := creatdb(nodeID) //创建数据库文件
		// fmt.Println("节点已启动")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close() // 在main函数中确保在所有操作完成后关闭数据库

		go p.tcpListen() //启动节点，这段这个完成返回的消息为tcpDial([]byte(info), p.messagePool[c.Digest].ClientAddr)

	} else {
		log.Fatal("无此节点编号！")
	}

	startTime = time.Now() // 获取当前时间
	fmt.Printf("当前时间：%v\n", startTime)

	select {
	case <-chanStop: // 主程序是否终止
		time.Sleep(time.Second)
		break
	}

}
