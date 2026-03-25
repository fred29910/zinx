package main

import (
	"time"

	"github.com/aceld/zinx/v3/znet"
)

func main() {
	s := znet.NewServer()

	// Start heartbeating detection.
	s.StartHeartBeat(5 * time.Second)

	s.Serve()
}
