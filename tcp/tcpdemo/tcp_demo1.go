package tcpdemo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	ListenerAddress = "127.0.0.1:10010"
)

// server
func TCPServer() {
	// 1. 开启监听服务
	listener, err := net.Listen("tcp", ListenerAddress)
	if err != nil {
		log.Fatalf("listener error: %s\n", err)
	}
	for {
		// 2.循环获取连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %s\n", err)
			continue
		}

		go serverProcess(conn) // 3. 启动协程处理连接.
	}
}

func serverProcess(conn net.Conn) {
	defer conn.Close()
	for {
		// 接收消息
		reader := bufio.NewReader(conn)
		var buffer [1024]byte
		n, err := reader.Read(buffer[:])
		if err != nil {
			if err == io.EOF { // 读取结束正常退出
				return
			}
			fmt.Printf("read from conn error: %v\n", err)
			return
		}
		revc := string(buffer[:n])
		fmt.Println("[Server Recvied]: ", string(revc))

		// 发送echo 给客户端
		echoMsg := "hello client"
		conn.Write([]byte(echoMsg))
	}
}

func TCPClient() {
	// 连接服务端
	conn, err := net.Dial("tcp", ListenerAddress)
	if err != nil {
		log.Fatalf("connect server %s error: %v\n", ListenerAddress, err)
	}

	// 通过输入进行发送消息 和 退出操作
	// q 退出
	// \n 遇到换行符发送给服务端
	inputReader := bufio.NewReader(os.Stdin)
	for {
		s, _ := inputReader.ReadString('\n') // 等待输入，直到出现 \n 换行结束
		s = strings.TrimSpace(s)             // 去空格

		// Q 退出
		if strings.ToUpper(s) == "Q" {
			return
		}
		// 给服务端发送消息
		_, err := conn.Write([]byte(s))
		if err != nil {
			log.Fatalf("send server error: %v\n", err)
		}
		// 接收服务端回复数据
		var buffer [1024]byte
		n, err := conn.Read(buffer[:])
		if err != nil {
			log.Fatalf("read server  data error: %v\n", err)
		}
		fmt.Println("[Client Recvied]: ", string(buffer[:n]))
	}
}
