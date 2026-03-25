package c_router

import (
	"fmt"
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloRouter) Handle(request ziface.IRequest) {
	slog.Debug("Call HelloZinxRouter Handle")

	slog.Debug(fmt.Sprint("recv from server : msgId=", request.GetMsgID()), ", data=", string(request.GetData()))
}
