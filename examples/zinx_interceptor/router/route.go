package router

import (
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

func (hr *HelloRouter) Handle(request ziface.IRequest) {
	slog.Info(string(request.GetData()))
}
