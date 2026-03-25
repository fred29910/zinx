package main

import (
	"log/slog"

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

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	s.AddRouter(1, &router.LoginRouter{})

	s.Serve()
}
