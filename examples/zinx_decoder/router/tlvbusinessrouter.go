package router

import (
	"fmt"
	"log/slog"

	"github.com/aceld/zinx/v3/zdecoder"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type TLVBusinessRouter struct {
	znet.BaseRouter
}

func (this *TLVBusinessRouter) Handle(request ziface.IRequest) {

	msgID := request.GetMessage().GetMsgID()
	slog.Debug(fmt.Sprintf("Call TLVRouter Handle %d %+v\n", msgID, request.GetMessage().GetData()))

	resp := request.GetResponse()
	if resp == nil {
		return
	}

	tlvData := resp.(zdecoder.TLVDecoder)
	slog.Debug(fmt.Sprintf("do msgid=0x00000001 data business %+v\n", tlvData))
}
