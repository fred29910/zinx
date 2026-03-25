package main

import (
	"fmt"
	"log/slog"

	"os"
	"os/signal"
	"time"

	"github.com/aceld/zinx/v3/zconf"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

const (
	PingType = 1
	PongType = 2
)

// releaseContextMiddleware releases the pooled Context after the middleware chain is done.
// (在整条 handler 链执行完成后归还 Context 对象池)
func releaseContextMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		defer c.Release()
		c.Next()
	}
}

// Hash 工作模式下，需要等待接受到client1的pong后，才会收到client2和client3的pong
// DynamicBind工作模式下，client2, client3 都会立马收到pong, 但client1的pong会被阻塞十秒后才收到
func PongHandler(client string) ziface.HandlerFunc {
	return func(c *ziface.Context) {
		// read server pong data
		slog.Info(fmt.Sprintf("---------client:%s, recv from server:%s, msgId=%d, data=%s ----------\n",
			client, c.Conn.RemoteAddr(), c.MsgID, string(c.Data)))
	}
}

func onClient1Start(conn ziface.IConnection) {
	slog.Info(fmt.Sprintf("client1 connection start, %s->%s\n", conn.LocalAddrString(), conn.RemoteAddrString()))
	//send ping
	err := conn.SendMsg(PingType, []byte("Ping From Client1"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

func onClient2Start(conn ziface.IConnection) {
	slog.Info(fmt.Sprintf("client2 connection start, %s->%s\n", conn.LocalAddrString(), conn.RemoteAddrString()))
	//send ping
	err := conn.SendMsg(PingType, []byte("Ping From Client2"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

func onClient3Start(conn ziface.IConnection) {
	slog.Info(fmt.Sprintf("client3 connection start, %s->%s\n", conn.LocalAddrString(), conn.RemoteAddrString()))
	//send ping
	err := conn.SendMsg(PingType, []byte("Ping From Client3"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

func main() {
	// Enable v3 context-based routing.
	zconf.GlobalObject.RouterSlicesMode = true

	//Create a client client
	client1 := znet.NewClient("127.0.0.1", 8999)
	client1.SetOnConnStart(onClient1Start)
	client1.GetMsgHandler().UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)
	client1.GetMsgHandler().AddRouterSlicesContext(PongType, PongHandler("client1"))
	client1.Start()

	time.Sleep(time.Second)

	client2 := znet.NewClient("127.0.0.1", 8999)
	client2.SetOnConnStart(onClient2Start)
	client2.GetMsgHandler().UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)
	client2.GetMsgHandler().AddRouterSlicesContext(PongType, PongHandler("client2"))
	client2.Start()

	time.Sleep(time.Second)

	client3 := znet.NewClient("127.0.0.1", 8999)
	client3.SetOnConnStart(onClient3Start)
	client3.GetMsgHandler().UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)
	client3.GetMsgHandler().AddRouterSlicesContext(PongType, PongHandler("client3"))
	client3.Start()

	//Prevent the process from exiting, waiting for an interrupt signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	client1.Stop()
	client2.Stop()
	client3.Stop()

	time.Sleep(time.Second)
}
