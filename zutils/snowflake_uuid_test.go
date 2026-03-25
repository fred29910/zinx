package zutils

import (
	"log/slog"
	"testing"
)

func TestSnowFlakeUUID(t *testing.T) {
	worker, err := NewIDWorker(1)
	if err != nil {
		slog.Debug("error occurred", "err", err)
		return
	}

	id, err := worker.NextID()
	if err != nil {
		slog.Debug("error occurred", "err", err)
		return
	}

	slog.Debug("ID:", "value", id)
}
