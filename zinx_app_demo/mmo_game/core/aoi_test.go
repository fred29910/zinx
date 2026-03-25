package core

import (
	"log/slog"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(100, 300, 4, 200, 450, 5)
	slog.Debug("aoiMgr info", "aoiMgr", aoiMgr)
}

func TestAOIManagerSuroundGrIDsByGID(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	for k := range aoiMgr.grIDs {
		// Get the surrounding nine grids of the current grid
		// (得到当前格子周边的九宫格)
		grIDs := aoiMgr.GetSurroundGrIDsByGID(k)

		// Get all IDs of the surrounding nine grids
		// (得到九宫格所有的IDs)
		slog.Debug("gID info", "gID", k, "grIDs len", len(grIDs))
		gIDs := make([]int, 0, len(grIDs))
		for _, grID := range grIDs {
			gIDs = append(gIDs, grID.GID)
		}
		slog.Debug("grID info", "grID", k, "gIDs", gIDs)
	}
}
