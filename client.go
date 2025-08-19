package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// // 用于客户端生成请求内容（.txt）
func readFile(filepath string) (data string, err error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil { //读取有误
		fmt.Println("读取出错，错误为:", err)
	}
	//如果读取成功，将内容显示在终端即可：
	// fmt.Printf("%v\n", string(content))
	fmt.Println("读取" + filepath + "成功！")
	return string(content), err
}

func countChanNum() {
	Mapmark := make(map[string]int)
	nCount := 0
	nFailCount := 0
	//fmt.Printf("countChanNum启动\n")
	for {
		fmt.Println("准备接受信息了")
		strB, ok := <-ChanNum
		fmt.Println("接受到的信息是：", strB)
		//	fmt.Printf("获取管道信息\n" + strB)
		if !ok {
			break
		}
		if strB == "failget" {
			nFailCount++
			str := "失败接收区块数" + strconv.Itoa(nFailCount)
			sendMSGToPrimary(str, nCount+nFailCount)
			if nCount+nFailCount == stopNum {
				bStopMark = true
				chanStop <- bStopMark
				return
			}
		}
		if _, ok := Mapmark[strB]; ok {
			Mapmark[strB]++
			fmt.Printf("我是多少%d, %d\n", Mapmark[strB], nodeCount)
			//if Mapmark[strB] == nodeCount/3+1 {
			if Mapmark[strB] == nodeCount/3+1 {
				//fmt.Println(strB)
				nCount++
				str := "成功接收区块数" + strconv.Itoa(nCount)
				fmt.Printf(str + "\n")
				//randSleep()
				time.Sleep(time.Millisecond * timeResentWait)
				sendMSGToPrimary(str, nCount)
				if nCount+nFailCount == stopNum {
					bStopMark = true
					chanStop <- bStopMark
					return
				}
			}
		} else {
			Mapmark[strB] = 1
			//fmt.Printf("我是多少%d\n", Mapmark[strB])
		}
	}

}

// 给主节点发送信息
func sendMSGToPrimary(data string, nID int) {
	r := new(Request) // Request是一个自定义结构体，分别保存下面赋值内容的信息
	r.Timestamp = time.Now().UnixNano()
	r.ClientAddr = clientAddr // "127.0.0.1:8888"
	r.Message.ID = nID

	//消息内容就是用户的输入
	r.Message.Content = strings.TrimSpace(data)
	br, err := json.Marshal(r) // 对r用json进行编码
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(string(br))               // 输出编码信息{"Content":"renyongwangshigedabendan","ID":4687201663,"Timestamp":1622769567507361000,"ClientAddr":"127.0.0.1:8888"}
	content := jointMessage(cRequest, br) // 合成请求信息
	//默认N0为主节点，直接把请求信息发送至N0
	tcpDial(content, nodeTable["N0"])
}

// 随机休眠1-10ms
/*func randSleep() {
	num := mrand.Intn(sleepTime) + 1
	for i := 0; i < num; i++ {
		time.Sleep(time.Millisecond)
	}
}*/

func clientSendMessageAndListen() {
	//开启客户端的本地监听（主要用来接收节点的reply信息）

	go clientTcpListen() //
	go countChanNum()

	fmt.Printf("客户端开启监听，地址：%s\n", clientAddr)
	fmt.Println(" ---------------------------------------------------------------------------------")
	fmt.Println("|  已进入PBFT测试Demo客户端，请启动全部节点后再发送消息！ :)  |")
	fmt.Println(" ---------------------------------------------------------------------------------")
	fmt.Println("请在下方输入要存入节点的信息：")

	// onechangenum, err := readFile()
	// if err != nil {
	// 	fmt.Println("Error reading from stdin")
	// 	panic(err)
	// }
	// go func() {
	// 	// onechangenum := "This is an example file.It's 11.49,February 1st 2024.I'm trying to make a PBFT demo with reed-solon algorithm."
	// 	num := 0
	// 	for {
	// 		num++

	// 		if num == 10 {
	// 			sendMSGToPrimary(onechangenum, getRandom())
	// 			break
	// 		}
	// 	}
	// }()

	go func() {
		// 创建一个定时器，每隔5秒触发一次
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// 启动一个goroutine来处理发送请求
		requestfiles := 12
		for i := 0; i < requestfiles; i++ {
			select {
			case <-ticker.C:
				// 从文件中读取内容
				filename := fmt.Sprintf("./clientfile/req%d.txt", i)
				onechangenum, err := readFile(filename)
				if err != nil {
					fmt.Println("Error reading from file:", err)
					continue // 如果读取失败，继续下一次循环
				}
				// 发送请求
				sendMSGToPrimary(onechangenum, getRandom())
			}
		}
		// 发送完请求后，退出函数
		fmt.Println("已发送%d次请求，结束共识。", requestfiles)
	}()
}

// 返回一个十位数的随机数，作为msgid
func getRandom() int {
	x := big.NewInt(10000000000)
	for {
		result, err := rand.Int(rand.Reader, x)
		if err != nil {
			log.Panic(err)
		}
		if result.Int64() > 1000000000 {
			return int(result.Int64())
		}
	}
}
