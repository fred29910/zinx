/**
* @Author: Aceld
* @Date: 2019/5/5 10:14
* @Mail: danbing.at@gmail.com
*
* 针对timer.go做单元测试，主要测试定时器相关接口 依赖模块delayFunc.go
 */
package ztimer

import (
	"log/slog"
	"testing"
	"time"
)

// 定义一个超时函数
func myFunc(v ...interface{}) {
	slog.Debug("function called", "No", v[0].(int), "delay_sec", v[1].(int))
}

func TestTimer(t *testing.T) {

	for i := 0; i < 5; i++ {
		go func(i int) {
			NewTimerAfter(NewDelayFunc(myFunc, []interface{}{i, 2 * i}), time.Duration(2*i)*time.Second).Run()
		}(i)
	}

	//主进程等待其他go，由于Run()方法是用一个新的go承载延迟方法，这里不能用waitGroup
	time.Sleep(1 * time.Minute)
}
