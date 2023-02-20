package cli

import (
	"context"
	"flag"
	"fmt"
	"go/scanner"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/config"
	"github.com/eugenenosenko/gopoly/generator"
	"github.com/eugenenosenko/gopoly/internal/xfs"
	"github.com/eugenenosenko/gopoly/poly"
	"github.com/eugenenosenko/gopoly/source"
)

type RecoveryFunc func(*App)

type App struct {
	Logf           func(format string, args ...any)
	Ctx            context.Context
	RecoveryFunc   RecoveryFunc
	ConfigProvider func() (*config.Config, error)
}

func (a *App) Exit(status int) {
	os.Exit(status)
}

func (a *App) Execute(info *RunInfo) {
	defer a.RecoveryFunc(a)
	if len(os.Args) > 2 && isInitCmd(os.Args[1]) {
		var code int
		if err := a.RunInit(); err != nil {
			log.Printf("Failed to execute init command %v", err)
			code = 3
		}
		a.Exit(code)
	}

	a.Logf("Run info: %s", info)

	sig := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(a.Ctx)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	defer func() {
		signal.Stop(sig)
		cancel()
	}()

	// trap cancellation signal on the context
	go func() {
		select {
		case <-sig:
			cancel()
		case <-ctx.Done():
		}
	}()
	a.execute()
}

func isInitCmd(arg string) bool {
	return arg == "init"
}

func (a *App) execute() {
	s, err := source.NewLoader(&source.Config{Logf: log.Printf, LoadFunc: source.LoadFromPackage})
	if err != nil {
		reportAndExit(a, errors.Wrapf(err, "creating source.Loader"))
	}

	gen, err := generator.NewTemplateGenerator(&generator.Config{
		Logf:     log.Printf,
		Provider: xfs.FileWriterProviderFunc(os.Create),
	})
	if err != nil {
		reportAndExit(a, errors.Wrapf(err, "creating codegen.Generator"))
	}

	runner, err := poly.NewClient(&poly.Config{
		Logf:          log.Printf,
		SourceLoader:  s,
		CodeGenerator: gen,
	})
	if err != nil {
		reportAndExit(a, errors.Wrapf(err, "creating poly.Client"))
	}

	conf, err := a.ConfigProvider()
	if err != nil {
		reportAndExit(a, errors.Wrapf(err, "getting types configuration"))
	}
	if err = runner.Run(a.Ctx, conf); err != nil {
		reportAndExit(a, err)
	}
	os.Exit(0)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: gopoly [-p package-path] [-c config-file] [-d decoding-strategy]"+
		" [-o output-file] [-m marker-method] [-t type-info]\n")
	flag.PrintDefaults()
}

func PrintTraceAndExit(app *App) {
	if r := recover(); r != nil {
		_, _ = fmt.Fprintf(os.Stderr, "panic occurred; %v, stack %s", r, debug.Stack())
	}
	app.Exit(3)
}

func reportAndExit(app *App, err error) {
	scanner.PrintError(os.Stderr, err)
	app.Exit(2)
}
