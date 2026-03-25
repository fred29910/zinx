/*
	Package zasync_op
	@Author：14March
	@File：async_worker.go
*/

package zasync_op

import (
	"fmt"
	"log/slog"
)

type AsyncWorker struct {
	taskQ chan func()
}

func (aw *AsyncWorker) process(asyncOp func()) {
	if asyncOp == nil {
		slog.Error(fmt.Sprintf("%v", "Async operation is empty."))
		return
	}

	if aw.taskQ == nil {
		slog.Error(fmt.Sprintf("%v", "Task queue has not been initialized."))
		return
	}

	aw.taskQ <- func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error(fmt.Sprintf("async process panic: %v", err))
			}
		}()

		// Execute async operation.(执行异步操作)
		asyncOp()
	}
}

func (aw *AsyncWorker) loopExecTask() {
	if aw.taskQ == nil {
		slog.Error(fmt.Sprintf("%v", "The task queue has not been initialized."))
		return
	}

	for {
		task := <-aw.taskQ
		if task != nil {
			task()
		}
	}
}
