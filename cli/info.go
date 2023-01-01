package cli

import (
	"fmt"
)

type RunInfo struct {
	ModSum     string
	ModVersion string
}

func (r *RunInfo) String() string {
	return fmt.Sprintf("mod-sum: %s mod-version: %s", r.ModSum, r.ModVersion)
}

var _ fmt.Stringer = (*RunInfo)(nil)
