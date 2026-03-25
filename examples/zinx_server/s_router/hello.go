package s_router

import (
	"fmt"
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type HelloZinxRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	slog.Debug("Call HelloZinxRouter Handle")
	// Read the data from the client first, then send back "ping...ping...ping"
	slog.Debug(fmt.Sprintf("recv from client : msgId=%d, data=%+v, len=%d", request.GetMsgID(), string(request.GetData()), len(request.GetData())))

	err := request.GetConnection().SendBuffMsg(3, []byte("Hello Zinx Router[FromServer]"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}
