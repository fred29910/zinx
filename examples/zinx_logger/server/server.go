package main

import (
	"log/slog"
	"time"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type TestRouter struct {
	znet.BaseRouter
}

// PreHandle -
func (t *TestRouter) PreHandle(req ziface.IRequest) {
	start := time.Now()

	slog.Debug("--> Call PreHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		slog.Debug("error occurred", "err", err)
	}
	elapsed := time.Since(start)
	slog.Debug("elapsed", "elapsed", elapsed)
}

// Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	slog.Debug("--> Call Handle")

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		slog.Debug("error occurred", "err", err)
	}
}

// PostHandle -
func (t *TestRouter) PostHandle(req ziface.IRequest) {
	slog.Debug("--> Call PostHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		slog.Debug("error occurred", "err", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &TestRouter{})
	// Note: Custom logger injection removed (zlog.SetLogger is no longer available).
	// Use slog.SetDefault() to configure a custom slog handler if needed.
	s.Serve()
}
