package main

import (
	"log/slog"
	"time"

	"github.com/aceld/zinx/v3/zconf"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// User-defined heartbeat message processing method
// 用户自定义的心跳检测消息处理方法
func myHeartBeatMsg(conn ziface.IConnection) []byte {
	return []byte("heartbeat, I am server, I am alive")
}

// User-defined handling method for remote connection not alive.
// 用户自定义的远程连接不存活时的处理方法
func myOnRemoteNotAlive(conn ziface.IConnection) {
	slog.Debug("myOnRemoteNotAlive is Called", "connID", conn.GetConnID(), "remoteAddr", conn.RemoteAddr())
	//关闭连接
	conn.Stop()
}

// User-defined heartbeat message handling function (用户自定义的心跳检测消息处理函数)
func myHeartBeatHandler(request ziface.IRequest) {
	slog.Debug("in myHeartBeatHandler", "msgId", request.GetMsgID(), "data", string(request.GetData()))
}

func main() {
	// Enable v3 RouterSlices mode so heartbeat can use RouterSlices.
	zconf.GlobalObject.RouterSlicesMode = true

	s := znet.NewServer()

	myHeartBeatMsgID := 88888

	// Start heartbeating detection. (启动心跳检测)
	s.StartHeartBeatWithOption(1*time.Second, &ziface.HeartBeatOption{
		MakeMsg:          myHeartBeatMsg,
		OnRemoteNotAlive: myOnRemoteNotAlive,
		RouterSlices:     []ziface.RouterHandler{myHeartBeatHandler},
		HeartBeatMsgID:   uint32(myHeartBeatMsgID),
	})

	s.Serve()
}
