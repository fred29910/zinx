package main

import (
	"errors"
	"log/slog"
	"time"

	"github.com/aceld/zinx/v3/zconf"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type TestRouter struct {
	znet.BaseRouter
}

var dealTimes = 0

// PreHandle -
func (t *TestRouter) PreHandle(req ziface.IRequest) {
	start := time.Now()

	slog.Debug("--> Call PreHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		slog.Debug("error occurred", "err", err)
	}
	elapsed := time.Since(start)
	slog.Debug("cost time", "elapsed", elapsed)
}

// Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	slog.Debug("--> Call Handle")

	if err := Err(); err != nil {
		req.Abort()
		slog.Debug("Insufficient permission")
	}

	dealTimes++
	req.GetConnection().AddCloseCallback(nil, nil, func() {
		slog.Debug("run close callback")
	})

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		slog.Debug("error occurred", "err", err)
	}

	if dealTimes == 5 {
		req.GetConnection().Stop()
	}

	time.Sleep(1 * time.Millisecond)
}

// PostHandle -
func (t *TestRouter) PostHandle(req ziface.IRequest) {
	slog.Debug("--> Call PostHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		slog.Debug("error occurred", "err", err)
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
		LogDir:             "./",
		LogFile:            "test.log",
		KcpFecDataShards:   10, //代表每10个原始数据块 发3个校验数据块
		KcpFecParityShards: 3,
	})
	s.AddRouter(1, &TestRouter{})
	s.SetOnConnStart(func(conn ziface.IConnection) {
		slog.Debug("--> OnConnStart")
	})
	s.SetOnConnStop(func(conn ziface.IConnection) {
		slog.Debug("--> OnConnStop")
	})
	s.Serve()
}
