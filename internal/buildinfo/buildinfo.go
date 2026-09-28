package buildinfo

import (
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/uvwt/agentdock/internal/executioncompat"
)

// This downstream version compares above the published 1.1.16102 and 1.1.17100
// distributions. The upstream source baseline is Workbench 1.1.8.
const Version = "1.1.18100"

// ProductName is the display identity from 1.1.7 onward; machine IDs remain stable.
const ProductName = "AgentDock Workbench"

var (
	Commit    string
	BuildDate string
)

type Info struct {
	ProductName            string `json:"product_name"`
	ExecutionPolicyVersion int    `json:"execution_policy_version"`
	Version                string `json:"version"`
	Commit                 string `json:"commit"`
	BuildDate              string `json:"build_date"`
	GoVersion              string `json:"go_version"`
	Platform               string `json:"platform"`
}

func Current() Info {
	info := Info{
		ProductName:            ProductName,
		ExecutionPolicyVersion: executioncompat.PolicyVersion,
		Version:                strings.TrimSpace(Version),
		Commit:                 strings.TrimSpace(Commit),
		BuildDate:              strings.TrimSpace(BuildDate),
		GoVersion:              runtime.Version(),
		Platform:               runtime.GOOS + "/" + runtime.GOARCH,
	}
	build, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range build.Settings {
			switch setting.Key {
			case "vcs.revision":
				if info.Commit == "" {
					info.Commit = strings.TrimSpace(setting.Value)
				}
			case "vcs.time":
				if info.BuildDate == "" {
					info.BuildDate = strings.TrimSpace(setting.Value)
				}
			}
		}
	}
	if len(info.Commit) > 12 {
		info.Commit = info.Commit[:12]
	}
	if info.Commit == "" {
		info.Commit = "unknown"
	}
	if info.BuildDate == "" {
		info.BuildDate = "unknown"
	}
	return info
}
