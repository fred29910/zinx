package main

import (
	"io"
	"log/slog"
	"net"

	"github.com/aceld/zinx/v3/zpack"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		slog.Debug("client start err, exit!", "err", err)
		return
	}

	dp := zpack.NewDataPack()
	msg, _ := dp.Pack(zpack.NewMsgPackage(1, []byte("async_op_router test=========>")))
	_, err = conn.Write(msg)
	if err != nil {
		slog.Debug("write error err ", "err", err)
		return
	}

	for {
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			slog.Debug("client read head err: ", "err", err)
			return
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			slog.Debug("client unpack head err: ", "err", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*zpack.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				slog.Debug("client unpack data err")
				return
			}

			slog.Debug("==> Client receive Msg", "ID", msg.ID, "len", msg.DataLen, "data", msg.Data)
		}
	}

}
