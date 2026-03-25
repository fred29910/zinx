package main

import (
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zinx_app_demo/mmo_game/pb"
	"github.com/aceld/zinx/v3/znet"
	"github.com/golang/protobuf/proto"
)

type PositionClientRouter struct {
	znet.BaseRouter
}

func (this *PositionClientRouter) Handle(request ziface.IRequest) {
	slog.Debug("Handle....")

	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		slog.Error("Position Unmarshal error", "err", err, "data", request.GetData())
		return
	}

	slog.Debug("recv from server", "msgId", request.GetMsgID(), "data", msg)
}

// 客户端自定义业务
func business(conn ziface.IConnection) {

	for {

		msg := &pb.Position{}
		msg.X = 1
		msg.Y = 2
		msg.Z = 3
		msg.V = 4

		data, err := proto.Marshal(msg)
		if err != nil {
			slog.Error("proto Marshal error", "err", err, "msg", msg)
			break
		}

		err = conn.SendMsg(0, data)
		if err != nil {
			slog.Debug("error occurred", "err", err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func DoClientConnectedBegin(conn ziface.IConnection) {
	conn.SetProperty("Name", "刘丹冰Aceld")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	slog.Debug("exit", "sig", sig)
}

func main() {
	client := znet.NewClient("127.0.0.1", 8999)

	client.SetOnConnStart(DoClientConnectedBegin)

	client.AddRouter(0, &PositionClientRouter{})

	client.Start()

	wait()
}
