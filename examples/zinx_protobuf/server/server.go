package main

import (
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zinx_app_demo/mmo_game/pb"
	"github.com/aceld/zinx/v3/znet"
	"github.com/golang/protobuf/proto"
)

type PositionServerRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PositionServerRouter) Handle(request ziface.IRequest) {

	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		slog.Error("Position Unmarshal error", "err", err, "data", request.GetData())
		return
	}

	slog.Debug("recv from client", "msgId", request.GetMsgID(), "data", msg)

	msg.X += 1
	msg.Y += 1
	msg.Z += 1
	msg.V += 1

	data, err := proto.Marshal(msg)
	if err != nil {
		slog.Error("proto Marshal error", "err", err, "msg", msg)
		return
	}

	err = request.GetConnection().SendMsg(0, data)

	if err != nil {
		slog.Error("error", "err", err)
	}
}

func main() {
	s := znet.NewServer()

	s.AddRouter(0, &PositionServerRouter{})

	s.Serve()
}
