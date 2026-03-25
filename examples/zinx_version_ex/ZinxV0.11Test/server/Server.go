/**
* @Author: Aceld
* @Date: 2019/4/30 17:42
* @Mail: danbing.at@gmail.com
*  ZinxV0.11测试，测试Zinx 日志模块功能 zlog模块
 */
package main

import (
	"fmt"
	"log/slog"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {

	slog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	slog.Debug(fmt.Sprint("recv from client : msgId=", request.GetMsgID()), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	slog.Debug("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	slog.Debug(fmt.Sprint("recv from client : msgId=", request.GetMsgID()), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello Zinx Router V0.10"))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

// 创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	slog.Debug("DoConnectionBegin is Called ... ")

	//设置两个连接属性，在连接创建之后
	slog.Debug(fmt.Sprint("Set conn Name, Home done!"))
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.jianshu.com/u/35261429b7f1")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		slog.Error("error", "err", err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		slog.Error("Conn Property Name = ", "detail", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		slog.Error("Conn Property Home = ", "detail", home)
	}

	slog.Debug("DoConneciotnLost is Called ... ")
}

func main() {
	//创建一个server句柄
	s := znet.NewServer()

	//注册连接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//开启服务
	s.Serve()
}
