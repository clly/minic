package main

import (
	"fmt"

	"go.clly.me/minic/sdk"
)

func main() {
	serviceInfo := sdk.VersionInfo()
	fmt.Printf("Main=%s\n", serviceInfo.Main)
	fmt.Printf("Go Version=%s\n", serviceInfo.GoVersion)
	fmt.Printf("VCS:\n")
	fmt.Printf("\tBranch=%s\n", serviceInfo.VCS.Branch)
	fmt.Printf("\tRevision=%s\n", serviceInfo.VCS.Revision)
	fmt.Printf("\tDirty=%v\n", serviceInfo.VCS.Dirty)
	fmt.Printf("\tCommitTime=%s\n", serviceInfo.VCS.CommitTime)
	fmt.Printf("\tTag=%s\n", serviceInfo.VCS.Tag)
	fmt.Printf("Build Info:\n")
	for k, v := range serviceInfo.Settings {
		fmt.Printf("\t%s=%s\n", k, v)
	}
}
