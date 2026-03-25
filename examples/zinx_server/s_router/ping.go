package s_router

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
	// Read the data from the client first, then send back "ping...ping...ping".
	slog.Debug(fmt.Sprintf("recv from client : msgId=%d, data=%+v, len=%d", request.GetMsgID(), string(request.GetData()), len(request.GetData())))

	err := request.GetConnection().SendMsg(2, []byte("pong-server"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}
