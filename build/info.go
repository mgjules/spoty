package build

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/JulesMike/spoty/logger"
	"go.uber.org/fx"
)

// Module exported for initialising a new build Info.
var Module = fx.Options(
	fx.Provide(New),
)

// Info contains the information about the build.
type Info struct {
	Revision   string    `json:"revision"`
	LastCommit time.Time `json:"last_commit"`
	DirtyBuild bool      `json:"dirty_build"`
}

// New returns a new instance of Info.
func New(logger *logger.Logger) (*Info, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, fmt.Errorf("failed to read build info")
	}

	info := Info{
		Revision:   "n/a",
		LastCommit: time.Time{},
		DirtyBuild: false,
	}

	for i := range bi.Settings {
		kv := &bi.Settings[i]

		switch kv.Key {
		case "vcs.revision":
			info.Revision = kv.Value
		case "vcs.time":
			hash, err := time.Parse(time.RFC3339, kv.Value)
			if err != nil {
				logger.Warnw("failed to parse vcs.time", "error", err, "value", kv.Value)
			}

			info.LastCommit = hash
		case "vcs.modified":
			info.DirtyBuild = kv.Value == "true"
		}
	}

	return &info, nil
}
