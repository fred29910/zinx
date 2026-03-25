package main

import (
	"log/slog"

	"github.com/aceld/zinx/v3/examples/zinx_interceptor/interceptors"
	"github.com/aceld/zinx/v3/examples/zinx_interceptor/router"
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

func main() {
	server := znet.NewServer()

	// Enable v3 context-based routing.
	zconf.GlobalObject.RouterSlicesMode = true

	// Register v3 context-based global middlewares.
	server.UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)

	// Register v3 context-based router handler.
	server.AddRouterSlicesContext(1, router.HelloRouter)

	// Add Custom Interceptor
	server.AddInterceptor(&interceptors.MyInterceptor{})

	slog.Info("zinx_interceptor server started")
	server.Serve()
}
