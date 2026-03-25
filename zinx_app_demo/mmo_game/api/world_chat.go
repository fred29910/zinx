package api

import (
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zinx_app_demo/mmo_game/core"
	"github.com/aceld/zinx/v3/zinx_app_demo/mmo_game/pb"
	"github.com/golang/protobuf/proto"
)

const playerContextKey = "player"

func RequirePlayer() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		pID, err := c.Conn.GetProperty("pID")
		if err != nil {
			slog.Debug("GetProperty pID error", "err", err)
			c.Conn.Stop()
			c.Abort()
			return
		}

		playerID, ok := pID.(int32)
		if !ok {
			slog.Debug("invalid player id type", "value", pID)
			c.Conn.Stop()
			c.Abort()
			return
		}

		player := core.WorldMgrObj.GetPlayerByPID(playerID)
		if player == nil {
			slog.Debug("player not found", "pID", playerID)
			c.Conn.Stop()
			c.Abort()
			return
		}

		c.Set(playerContextKey, player)
		c.Next()
	}
}

func currentPlayer(c *ziface.Context) (*core.Player, bool) {
	playerValue, exists := c.Get(playerContextKey)
	if !exists {
		slog.Debug("player missing in context")
		c.Abort()
		return nil, false
	}

	player, ok := playerValue.(*core.Player)
	if !ok || player == nil {
		slog.Debug("invalid player in context", "value", playerValue)
		c.Abort()
		return nil, false
	}

	return player, true
}

// WorldChat handles world chat messages in v3 context mode.
// (WorldChat 使用 v3 Context 模式处理世界聊天消息)
func WorldChat(c *ziface.Context) {
	msg := &pb.Talk{}
	if err := proto.Unmarshal(c.Data, msg); err != nil {
		slog.Debug("Talk Unmarshal error ", "err", err)
		return
	}

	player, ok := currentPlayer(c)
	if !ok {
		return
	}

	player.Talk(msg.Content)
}
