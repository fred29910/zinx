package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aceld/zinx/v3/zconf"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// 客户端自定义业务
func business(conn ziface.IConnection) {
	for {
		// Keep msgID consistent with server-side router slices.
		// (客户端发送 -> server 的 msgID=100 -> server 回包 msgID=2 -> 客户端处理)
		err := conn.SendMsg(100, []byte("ping ping ping ..."))
		if err != nil {
			slog.Debug("error occurred", "err", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// 创建连接的时候执行
func DoClientConnectedBegin(conn ziface.IConnection) {
	//设置两个连接属性，在连接创建之后
	conn.SetProperty("Name", "刘丹冰")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

// releaseContextMiddleware releases the pooled Context after the middleware chain is done.
// (在整条 handler 链执行完成后归还 Context 对象池)
func releaseContextMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		defer c.Release()
		c.Next()
	}
}

func PingRouter(c *ziface.Context) {
	slog.Debug("Call PingRouter Handle")
	slog.Debug("recv from server", "msgID", c.MsgID, "data", string(c.Data), "len", len(c.Data))

	if err := c.Conn.SendBuffMsg(1, []byte("Hello[from client]")); err != nil {
		slog.Error("error", "err", err)
	}
}

func HelloRouter(c *ziface.Context) {
	slog.Debug("Call HelloZinxRouter Handle")
	slog.Debug("recv from server", "msgID", c.MsgID, "data", string(c.Data), "len", len(c.Data))
}

func main() {
	// Enable v3 context-based router slices.
	zconf.GlobalObject.Mode = ""
	zconf.GlobalObject.RouterSlicesMode = true

	// Create a Client.
	client := znet.NewWsClient("127.0.0.1", 9000)

	// Register v3 context-based global middlewares (for the client's message processing pipeline).
	client.GetMsgHandler().UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)

	// Register business logic via context-based router slices.
	client.GetMsgHandler().AddRouterSlicesContext(2, PingRouter)   // server -> msgID=2
	client.GetMsgHandler().AddRouterSlicesContext(3, HelloRouter)  // server -> msgID=3

	// Add business logic for when the connection is first established.(添加首次建立连接时的业务)
	client.SetOnConnStart(DoClientConnectedBegin)

	// Start the client.
	client.Start()

	// Wait for either client error or process signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := <-client.GetErrChan(); err != nil {
			slog.Error("client err", "err", err)
			// Unblock main by emitting a signal into the same channel.
			// (避免在 goroutine 中再去读取 `c`，否则会抢走 main 的信号)
			c <- os.Interrupt
		}
	}()

	sig := <-c
	slog.Debug("exit", "sig", sig)
	// Clean up the client.(清理客户端)
	client.Stop()
	time.Sleep(2 * time.Second)
}
