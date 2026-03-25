package api

import (
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zinx_app_demo/mmo_game/pb"
	"github.com/golang/protobuf/proto"
)

// Move handles movement messages in v3 context mode.
// (Move 使用 v3 Context 模式处理玩家移动消息)
func Move(c *ziface.Context) {
	msg := &pb.Position{}
	if err := proto.Unmarshal(c.Data, msg); err != nil {
		slog.Debug("Move: Position Unmarshal error ", "err", err)
		return
	}

	player, ok := currentPlayer(c)
	if !ok {
		return
	}

	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
