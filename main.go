package main

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/eugenenosenko/gopoly/cli"
)

func main() {
	app := cli.App{
		Logf:           log.Printf,
		RecoveryFunc:   cli.PrintTraceAndExit,
		Ctx:            context.Background(),
		ConfigProvider: cli.MergedConfigProvider,
	}
	app.Execute(newRunInfo())
}

func newRunInfo() *cli.RunInfo {
	var version, sum string
	if info, exists := debug.ReadBuildInfo(); exists {
		version = info.Main.Version
		sum = info.Main.Sum
	}
	return &cli.RunInfo{ModSum: sum, ModVersion: version}
}
