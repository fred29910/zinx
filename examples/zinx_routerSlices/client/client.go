package main

import (
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
	msg, _ := dp.Pack(zpack.NewMsgPackage(1, []byte("ZinxPing")))
	_, err = conn.Write(msg)
	if err != nil {
		slog.Debug("write error err ", "err", err)
		return
	}

}
