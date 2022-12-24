package poly

import (
	"fmt"
	"runtime/debug"
)

type runInfo struct {
	modSum     string
	modVersion string
}

func (r *runInfo) String() string {
	return fmt.Sprintf("mod-sum: %s mod-version: %s", r.modSum, r.modVersion)
}

func newRunInfo() *runInfo {
	var version, sum string
	if info, exists := debug.ReadBuildInfo(); exists {
		version = info.Main.Version
		sum = info.Main.Sum
	}
	return &runInfo{modSum: sum, modVersion: version}
}

var _ fmt.Stringer = (*runInfo)(nil)
