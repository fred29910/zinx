package main

import (
	"errors"
	"log/slog"
	"time"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

type TestRouter struct {
	znet.BaseRouter
}

// PreHandle -
func (t *TestRouter) PreHandle(req ziface.IRequest) {
	start := time.Now()

	slog.Debug("--> Call PreHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		slog.Debug("error occurred", "err", err)
	}
	elapsed := time.Since(start)
	slog.Debug("cost time", "elapsed", elapsed)
}

// Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	slog.Debug("--> Call Handle")

	// Simulated scenario - In the event of an expected error such as incorrect permissions or incorrect information,
	// subsequent function execution will be stopped, but this function will be fully executed.
	// 模拟场景- 出现意料之中的错误 如权限不对或者信息错误 则停止后续函数执行，但是次函数会执行完毕
	if err := Err(); err != nil {
		req.Abort()
		slog.Debug("Insufficient permission")
	}

	// Simulation scenario - In case of a certain situation, repeat the above operation.
	// 模拟场景- 出现某种情况，重复上面的操作
	/*
		if err := Err(); err != nil {
			req.Goto(znet.PRE_HANDLE)
			slog.Debug("repeat")
		}
	*/

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		slog.Debug("error occurred", "err", err)
	}

	time.Sleep(1 * time.Millisecond)
}

// PostHandle -
func (t *TestRouter) PostHandle(req ziface.IRequest) {
	slog.Debug("--> Call PostHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		slog.Debug("error occurred", "err", err)
	}
}

func Err() error {
	//Specific Business Operation (具体业务操作)
	return errors.New("Test")
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &TestRouter{})
	s.Serve()
}
