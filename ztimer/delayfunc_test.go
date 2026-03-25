/**
* @Author: Aceld
* @Date: 2019/4/30 15:17
* @Mail: danbing.at@gmail.com
*
*  针对 delayFunc.go 做单元测试，主要测试延迟函数结构体是否正常使用
 */
package ztimer

import (
	"log/slog"
	"testing"
)

func SayHello(message ...interface{}) {
	slog.Debug("delay message", "msg0", message[0].(string), "msg1", message[1].(string))
}

func TestDelayfunc(t *testing.T) {
	df := NewDelayFunc(SayHello, []interface{}{"hello", "zinx!"})
	slog.Debug("df.String()", "value", df.String())
	df.Call()
}
