package api

import (
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zinx_app_demo/mmo_game/core"
	"github.com/aceld/zinx/v3/zinx_app_demo/mmo_game/pb"
	"github.com/aceld/zinx/v3/znet"
	"github.com/golang/protobuf/proto"
)

// MoveApi Player movement
// MoveApi 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (*MoveApi) Handle(request ziface.IRequest) {
	//1. MoveApi Player movement
	// (1. 将客户端传来的proto协议解码)
	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		slog.Debug("Move: Position Unmarshal error ", "err", err)
		return
	}

	// 2. Identify which player sent the current message, retrieve from the connection property pID
	// (2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取)
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		slog.Debug("GetProperty pID error", "err", err)
		request.GetConnection().Stop()
		return
	}

	// slog.Debug("user move", "pID", pID, "X", msg.X, "Y", msg.Y, "Z", msg.Z, "V", msg.V)

	// 3. Get the player object based on pID
	// (3. 根据pID得到player对象)
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	// 4. Have the player object initiate the broadcast of movement position information
	// (4. 让player对象发起移动位置信息广播)
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
