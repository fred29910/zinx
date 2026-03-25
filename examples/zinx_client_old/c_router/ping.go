package c_router

import (
	"fmt"
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	slog.Debug("Call PingRouter Handle")

	slog.Debug(fmt.Sprint("recv from server : msgId=", request.GetMsgID()), ", data=", string(request.GetData()))

	if err := request.GetConnection().SendBuffMsg(1, []byte("Hello[from client]")); err != nil {
		slog.Error("error", "err", err)
	}
}
