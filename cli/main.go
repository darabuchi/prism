package main

import (
	"os"

	"github.com/prism/cli/cmd"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// 设置版本信息
	cmd.SetVersionInfo(version, buildTime, gitCommit)
	
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}