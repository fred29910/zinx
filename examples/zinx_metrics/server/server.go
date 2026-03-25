package main

import (
	"fmt"
	"log/slog"

	"github.com/aceld/zinx/v3/examples/zinx_server/s_router"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
	slog.Info("DoConnectionBegin is Called ...")

	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	if name, err := conn.GetProperty("Name"); err == nil {
		slog.Info(fmt.Sprintf("Conn Property Name = %v", name))
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		slog.Info(fmt.Sprintf("Conn Property Home = %v", home))
	}

	slog.Info("Conn is Lost")
}

// usage:$  curl 0.0.0.0:20004/metrics
// to get Metrics
func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	s.Serve()
}
