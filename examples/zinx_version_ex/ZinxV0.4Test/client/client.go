package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

/*
模拟客户端
*/
func main() {

	slog.Debug("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		slog.Debug("client start err, exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("Zinx V0.3"))
		if err != nil {
			slog.Debug("write error err ", "err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			slog.Debug("read buf error ")
			return
		}

		fmt.Printf(" server call back : %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
}
