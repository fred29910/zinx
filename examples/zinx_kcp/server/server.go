package main

import (
	"errors"
	"log/slog"
	"time"

	"github.com/aceld/zinx/v3/zconf"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

var dealTimes = 0

// releaseContextMiddleware releases the pooled Context after the middleware chain is done.
// (在整条 handler 链执行完成后归还 Context 对象池)
func releaseContextMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		defer c.Release()
		c.Next()
	}
}

func TestHandler(c *ziface.Context) {
	// ---- PreHandle (merged into handler for v3 slices mode) ----
	start := time.Now()

	slog.Debug("--> Call PreHandle")
	if err := c.Conn.SendMsg(0, []byte("test1")); err != nil {
		slog.Debug("error occurred", "err", err)
	}
	elapsed := time.Since(start)
	slog.Debug("cost time", "elapsed", elapsed)

	// ---- Handle ----
	slog.Debug("--> Call Handle")

	// In v1, req.Abort() will skip PostHandle; mimic that behavior here.
	skipPostHandle := false
	if err := Err(); err != nil {
		c.Abort()
		skipPostHandle = true
		slog.Debug("Insufficient permission")
	}

	dealTimes++

	c.Conn.AddCloseCallback(nil, nil, func() {
		slog.Debug("run close callback")
	})

	if err := c.Conn.SendMsg(0, []byte("test2")); err != nil {
		slog.Debug("error occurred", "err", err)
	}

	if dealTimes == 5 {
		c.Conn.Stop()
	}

	time.Sleep(1 * time.Millisecond)

	// ---- PostHandle (only when not aborted) ----
	if !skipPostHandle {
		slog.Debug("--> Call PostHandle")
		if err := c.Conn.SendMsg(0, []byte("test3")); err != nil {
			slog.Debug("error occurred", "err", err)
		}
	}
}

func Err() error {
	//Specific Business Operation (具体业务操作)
	return errors.New("Test")
}

func main() {
	s := znet.NewUserConfServer(&zconf.Config{
		Mode:               "kcp",
		KcpPort:            7777,
		KcpRecvWindow:      128,
		KcpSendWindow:      128,
		KcpStreamMode:      true,
		KcpACKNoDelay:      false,
		KcpFecDataShards:   10, //代表每10个原始数据块 发3个校验数据块
		KcpFecParityShards: 3,
	})

	// Enable v3 context-based routing.
	zconf.GlobalObject.RouterSlicesMode = true
	s.UseContext(
		releaseContextMiddleware(),
		znet.RecoveryMiddleware(),
		znet.SlogLoggerMiddleware(),
	)

	s.AddRouterSlicesContext(1, TestHandler)
	s.SetOnConnStart(func(conn ziface.IConnection) {
		slog.Debug("--> OnConnStart")
	})
	s.SetOnConnStop(func(conn ziface.IConnection) {
		slog.Debug("--> OnConnStop")
	})
	s.Serve()
}
