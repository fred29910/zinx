package main

import (
	"fmt"
	"log/slog"

	"sync/atomic"
	"time"

	"github.com/aceld/zinx/v3/zconf"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	slog.Debug(fmt.Sprintf("OnConnectionAdd: %v", conn.GetConnection().RemoteAddr()))
}

func OnConnectionLost(conn ziface.IConnection) {
	slog.Debug(fmt.Sprintf("OnConnectionLost: %v", conn.GetConnection().RemoteAddr()))
}

var Block = int32(1)

// releaseContextMiddleware releases the pooled Context after the middleware chain is done.
// (在整条 handler 链执行完成后归还 Context 对象池)
func releaseContextMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		defer c.Release()
		c.Next()
	}
}

// 模拟阻塞操作 (v3 RouterSlicesContext handler)
func blockHandler(c *ziface.Context) {
	// read client data
	slog.Info(fmt.Sprintf("recv from client:%s, msgId=%d, data=%s\n", c.Conn.RemoteAddr(), c.MsgID, string(c.Data)))

	// 第一次处理时，模拟任务阻塞操作, Hash 模式下，后面的连接的任务得不到处理
	// DynamicBind 模式下，看后面的连接的任务会得到即使处理，不会因为前面连接的任务阻塞而得不到处理
	// 这里只模拟一次阻塞操作。
	if atomic.CompareAndSwapInt32(&Block, 1, 0) {
		slog.Info(fmt.Sprintf("blockRouter handle start, msgId=%d, remote:%v\n", c.MsgID, c.Conn.RemoteAddr()))
		time.Sleep(time.Second * 10)
		//阻塞操作结束
		slog.Info(fmt.Sprintf("blockRouter handle end, msgId=%d, remote:%v\n", c.MsgID, c.Conn.RemoteAddr()))
	}

	err := c.Conn.SendMsg(2, []byte("pong from server"))
	if err != nil {
		slog.Error("error", "err", err)
		return
	}
	slog.Info(fmt.Sprintf("send pong over, client:%s\n", c.Conn.RemoteAddr()))
}

func main() {
	// Load server settings to keep the example self-contained.
	s := znet.NewUserConfServer(&zconf.Config{
		Name:              "zinx server DynamicBind Mode Demo",
		Host:              "127.0.0.1",
		TCPPort:           8999,
		MaxConn:           12000,
		WorkerPoolSize:    1,
		MaxWorkerTaskLen: 50,
		WorkerMode:        zconf.WorkerModeDynamicBind,
		RouterSlicesMode:  true,
	})

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// Register v3 context-based global middlewares.
	s.UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)

	s.AddRouterSlicesContext(1, blockHandler)

	s.Serve()
}
