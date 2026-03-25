package main

import (
	"log/slog"

	"github.com/aceld/zinx/v3/zconf"

	"github.com/aceld/zinx/v3/examples/zinx_async_op/router"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	slog.Debug("zinx_async_op OnConnectionAdd ===>")
}

func OnConnectionLost(conn ziface.IConnection) {
	slog.Debug("zinx_async_op OnConnectionLost ===>")
}

// releaseContextMiddleware releases the pooled Context after the middleware chain is done.
// (在整条 handler 链执行完成后归还 Context 对象池)
func releaseContextMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		defer c.Release()
		c.Next()
	}
}

func main() {
	// Enable v3 context-based routing.
	zconf.GlobalObject.RouterSlicesMode = true

	s := znet.NewServer()

	// Register v3 context-based global middlewares.
	s.UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	s.AddRouterSlicesContext(1, router.LoginHandler)

	s.Serve()
}
