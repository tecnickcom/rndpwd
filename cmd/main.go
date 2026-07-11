// Package main is an example web service to generate random passwords.
package main

import (
	"log/slog"
	"os"

	"github.com/tecnickcom/nurago/pkg/bootstrap"
	"github.com/tecnickcom/nurago/pkg/logsrv"
	"github.com/tecnickcom/nurago/pkg/logutil"
	"github.com/tecnickcom/rndpwd/internal/cli"
)

var (
	// programVersion contains the version of the application injected at compile time.
	programVersion = "0.0.0" //nolint:gochecknoglobals

	// programRelease contains the release of the application injected at compile time.
	programRelease = "0" //nolint:gochecknoglobals
)

// exitFn defines the exit function and can be overwritten for testing.
var exitFn = os.Exit //nolint:gochecknoglobals

func main() {
	// set default logger
	logattr := []logutil.Attr{
		slog.String("program", cli.AppName),
		slog.String("version", programVersion),
		slog.String("release", programRelease),
	}
	logcfg, _ := logutil.NewConfig(
		logutil.WithOutWriter(os.Stderr),
		logutil.WithFormat(logutil.FormatJSON),
		logutil.WithLevel(logutil.LevelDebug),
		logutil.WithCommonAttr(logattr...),
	)
	l := logsrv.NewLogger(logcfg)

	rootCmd, err := cli.New(programVersion, programRelease, bootstrap.Bootstrap)
	if err != nil {
		l.With(slog.Any("error", err)).Error("UNABLE TO START THE PROGRAM")
		exitFn(1)

		// exitFn normally terminates the process; guard against test doubles
		// that return, so a nil rootCmd is never executed.
		return
	}

	// execute the root command and log errors (if any)
	err = rootCmd.Execute()
	if err != nil {
		l.With(slog.Any("error", err)).Error("UNABLE TO RUN THE COMMAND")
		exitFn(2)
	}
}
