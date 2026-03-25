package router

import (
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
)

// HelloRouter is a v3 context-based handler.
func HelloRouter(c *ziface.Context) {
	slog.Info("HelloRouter recv", "msgID", c.MsgID, "data", string(c.Data), "len", len(c.Data))
}
