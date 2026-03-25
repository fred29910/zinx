package znet

import (
	"io"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zpack"
)

// run in terminal:
// go test -v ./znet -run=TestServer

/*
ClientTest client
*/
func ClientTest(i uint32) {

	slog.Debug("Client Test ... start")

	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		slog.Debug("client start err, exit!")
		return
	}

	for {
		dp := zpack.Factory().NewPack(ziface.ZinxDataPack)
		msg, _ := dp.Pack(zpack.NewMsgPackage(i, []byte("client test message")))
		_, err := conn.Write(msg)
		if err != nil {
			slog.Debug("client write err: ", "err", err)
			return
		}

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

		time.Sleep(time.Second)
	}
}

/*
	server
*/

type PingRouter struct {
	BaseRouter
}

// Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	slog.Debug("Call Router PreHandle")
	err := request.GetConnection().SendMsg(1, []byte("before ping ....\n"))
	if err != nil {
		slog.Debug("preHandle SendMsg err: ", "err", err)
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	slog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	slog.Debug("recv from client", "msgID", request.GetMsgID(), "data", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping\n"))
	if err != nil {
		slog.Debug("Handle SendMsg err: ", "err", err)
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	slog.Debug("Call Router PostHandle")
	err := request.GetConnection().SendMsg(1, []byte("After ping .....\n"))
	if err != nil {
		slog.Debug("Post SendMsg err: ", "err", err)
	}
}

type HelloRouter struct {
	BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	slog.Debug("call helloRouter Handle")
	slog.Debug("receive from client", "msgID", request.GetMsgID(), "data", string(request.GetData()))

	err := request.GetConnection().SendMsg(2, []byte("hello zix hello Router"))
	if err != nil {
		slog.Debug("error occurred", "err", err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	slog.Debug("DoConnectionBegin is Called ... ")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		slog.Debug("error occurred", "err", err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	slog.Debug("DoConnectionLost is Called ... ")
}

func TestServer(t *testing.T) {
	s := NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	go ClientTest(1)
	go ClientTest(2)

	go s.Serve()

	select {
	case <-time.After(time.Second * 10):
		return
	}
}

func TestServerDeadLock(t *testing.T) {
	s := NewServer()

	s.Start()
	time.Sleep(time.Second * 1)

	go func() {
		_, _ = net.Dial("tcp", "127.0.0.1:8999")
	}()
	time.Sleep(time.Second * 1)
	s.Stop()
}

type CloseConnectionBeforeSendMsgRouter struct {
	BaseRouter
}

type DemoPacket struct {
	zpack.DataPack
}

func (d *DemoPacket) Pack(msg ziface.IMessage) ([]byte, error) {
	time.Sleep(time.Second * 1)
	return d.DataPack.Pack(msg)
}

func (br *CloseConnectionBeforeSendMsgRouter) Handle(req ziface.IRequest) {
	connection := req.GetConnection()
	msg := "Zinx server response message for CloseConnectionBeforeSendMsgRouter"
	connection.Stop()
	_ = connection.SendMsg(1, []byte(msg))
	slog.Debug("send:", "msg", msg)
}

func TestCloseConnectionBeforeSendMsg(t *testing.T) {
	s := NewServer()
	s.AddRouter(1, &CloseConnectionBeforeSendMsgRouter{})

	s.Start()
	time.Sleep(time.Second * 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		conn, _ := net.Dial("tcp", "127.0.0.1:8999")
		dp := zpack.Factory().NewPack(ziface.ZinxDataPack)
		msg := "Zinx client request message for CloseConnectionBeforeSendMsgRouter"
		pack, _ := dp.Pack(zpack.NewMsgPackage(1, []byte(msg)))
		_, _ = conn.Write(pack)
		slog.Debug("send:", "msg", msg)
		buffer := make([]byte, 1024)
		readLen, _ := conn.Read(buffer)
		slog.Debug("received all data", "data", string(buffer[dp.GetHeadLen():readLen]))
		wg.Done()
	}()
	wg.Wait()
	s.Stop()
}
