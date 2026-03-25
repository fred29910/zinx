package main

import (
	"io"
	"log/slog"
	"time"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zpack"
	"github.com/xtaci/kcp-go/v5"
)

// 模拟客户端
func main() {
	slog.Debug("Client Test ... start")
	// Replace net.Dial with kcp.DialWithOptions
	conn, err := kcp.Dial("127.0.0.1:7777")
	if err != nil {
		slog.Debug("client start err, exit!")
		return
	}

	dp := zpack.Factory().NewPack(ziface.ZinxDataPack)
	sendMsg, _ := dp.Pack(zpack.NewMsgPackage(1, []byte("client test message")))
	_, err = conn.Write(sendMsg)
	if err != nil {
		slog.Debug("client write err: ", "err", err)
		return
	}

	for {
		// Read the "head" section from the stream first. (先读出流中的head部分)
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			slog.Debug("client read head err: ", "err", err)
			return
		}

		// Unpack the headData byte stream into msg. (将headData字节流 拆包到msg中)
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			slog.Debug("client unpack head err: ", "err", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// Read the "data" section from the stream. (再读出流中的data部分)
			msg := msgHead.(*zpack.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// read from io.Reader into msg.Data (根据dataLen从io中读取字节流)
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				slog.Debug("client unpack data err")
				return
			}

			slog.Debug("==> Client receive Msg", "ID", msg.ID, "len", msg.DataLen, "data", msg.Data)

			time.Sleep(1 * time.Second)
			_, err = conn.Write(sendMsg)
			if err != nil {
				slog.Debug("client write err: ", "err", err)
				return
			}
		}
	}
}
