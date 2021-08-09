package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func process(conn net.Conn) {
	// 处理完关闭连接
	defer conn.Close()

	ch := make(chan []byte, 1)

	times := 0

	// 针对当前连接做发送和接受操作
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte

		// block here until client send sth
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			break
		}

		ch <- buf[:n]

		go checkHeartBeat(ch, conn, times)

		recv := string(buf[:n])
		fmt.Printf("收到的数据：%v\n", recv)

		// 将接受到的数据返回给客户端
		_, err = conn.Write([]byte("ok"))
		if err != nil {
			fmt.Printf("write from conn failed, err:%v\n", err)
			break
		}

		times++
	}
}

func checkHeartBeat(ch chan []byte, conn net.Conn, times int) {
	for {
		select {
		case <-time.After(time.Second * 4):
			fmt.Println("time out, close conn")
			conn.Close()
		case he := <-ch:
			fmt.Println(times, " get heartbeat:", string(he))
		}
		return
	}
}

func main() {
	// 建立 tcp 服务
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("listen failed, err:%v\n", err)
		return
	}

	for {
		// 等待客户端建立连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}
		// 启动一个单独的 goroutine 去处理连接
		go process(conn)
	}
}
