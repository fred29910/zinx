package router

import (
	"fmt"
	"log/slog"

	"encoding/hex"

	"github.com/aceld/zinx/v3/zdecoder"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type HtlvCrcBusinessRouter struct {
	znet.BaseRouter
}

func (this *HtlvCrcBusinessRouter) Handle(request ziface.IRequest) {

	//MsgID
	msgID := request.GetMessage().GetMsgID()
	slog.Debug(fmt.Sprintf("Call HtlvCrcBusinessRouter Handle %d %s\n", msgID, hex.EncodeToString(request.GetMessage().GetData())))

	resp := request.GetResponse()
	if resp == nil {
		return
	}

	tlvData := resp.(zdecoder.HtlvCrcDecoder)

	slog.Debug(fmt.Sprintf("do msgid=0x10 data business %+v\n", tlvData))
}
