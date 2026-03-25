package main

import (
	"log/slog"

	"github.com/aceld/zinx/v3/zconf"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

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
	// Read the data from the client first, then send back "ping...ping...ping".
	slog.Debug("recv from client", "msgId", c.MsgID, "data", string(c.Data), "len", len(c.Data))

	err := c.Conn.SendMsg(2, []byte("pong-server"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

func HelloZinxRouter(c *ziface.Context) {
	slog.Debug("Call HelloZinxRouter Handle")
	// Read the data from the client first, then send back "ping...ping...ping"
	slog.Debug("recv from client", "msgId", c.MsgID, "data", string(c.Data), "len", len(c.Data))

	err := c.Conn.SendBuffMsg(3, []byte("Hello Zinx Router[FromServer]"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

func main() {
	// Set up as WebSocket before starting. (在启动之前设置为 websocket)
	zconf.GlobalObject.Mode = ""
	// Enable V3 Router Slices Mode
	zconf.GlobalObject.RouterSlicesMode = true

	s := znet.NewServer()

	// 注册全局中间件
	s.UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)

	// 改用新的 SlicesContext 处理方法注入逻辑
	s.AddRouterSlicesContext(100, PingRouter)
	s.AddRouterSlicesContext(1, HelloZinxRouter)

	s.Serve()
}
