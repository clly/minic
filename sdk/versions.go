package sdk

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

var serviceInfo *ServiceInfo

// These are expected to be set by LD flags
var branch = ""
var tag = ""

type VCS struct {
	Branch     string
	Tag        string
	Dirty      bool
	CommitTime time.Time
	Revision   string
}

type ServiceInfo struct {
	GoVersion string
	Deps      []*debug.Module
	Main      string
	Settings  map[string]string
	VCS       *VCS
}

func VersionInfo() *ServiceInfo {
	if serviceInfo != nil {
		return serviceInfo
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	vcs := &VCS{
		Branch: branch,
		Tag:    tag,
	}
	serviceInfo = &ServiceInfo{
		GoVersion: info.GoVersion,
		Deps:      info.Deps,
		Main:      info.Main.Path,
		VCS:       vcs,
	}

	settings := make(map[string]string)
	var err error
	for _, v := range info.Settings {
		switch v.Key {
		case "vcs.time":
			vcs.CommitTime, err = time.Parse(time.RFC3339, v.Value)
			if err != nil {
				settings[v.Key] = v.Value
				fmt.Fprintln(os.Stderr, "failed to parse vcs time", err)
			}
		case "vcs.modified":
			vcs.Dirty, err = strconv.ParseBool(v.Value)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to parse vcs.modified", err)
			}
		case "vcs.revision":
			vcs.Revision = v.Value
		default:
			settings[v.Key] = v.Value
		}
	}
	serviceInfo.Settings = settings
	return serviceInfo
}
