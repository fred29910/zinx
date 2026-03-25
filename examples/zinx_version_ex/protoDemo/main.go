package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/aceld/zinx/v3/examples/zinx_version_ex/protoDemo/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name:   "XiaoYuer",
		Age:    16,
		Emails: []string{"xiao_yu_er@sina.com", "yu_er@sina.cn"},
		Phones: []*pb.PhoneNumber{
			{
				Number: "13113111311",
				Type:   pb.PhoneType_MOBILE,
			},
			{
				Number: "14141444144",
				Type:   pb.PhoneType_HOME,
			},
			{
				Number: "19191919191",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	data, err := proto.Marshal(person)
	if err != nil {
		slog.Debug("marshal err:", "err", err)
	}

	fmt.Println(hex.EncodeToString(data))

	newdata := &pb.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		slog.Debug("unmarshal err:", "err", err)
	}
	fmt.Println(newdata)
}
