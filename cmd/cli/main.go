// Command cli is the interactive UI for the STALART JVM wrapper:
// install/uninstall IFEO hooks, pick configs. The IFEO debugger binary
// is service.exe in the same directory.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"stalart-wrapper/internal/installer"
	"stalart-wrapper/internal/logging"
	"stalart-wrapper/internal/ui"
)

func main() {
	closeLog, err := logging.Setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[log] %v\n", err)
	}
	defer closeLog()

	slog.Info("cli startup", "args_count", len(os.Args)-1)

	if handled, code := handleCLI(); handled {
		os.Exit(code)
	}

	if err := ui.Run(); err != nil {
		slog.Error("ui failed", "err", err)
		fmt.Fprintf(os.Stderr, "[ui] %v\n", err)
		os.Exit(1)
	}
}

func handleCLI() (handled bool, code int) {
	if len(os.Args) < 2 {
		return false, 0
	}
	switch os.Args[1] {
	case "--install":
		if err := installer.Install(); err != nil {
			slog.Error("install failed", "err", err)
			fmt.Fprintf(os.Stderr, "[install] %v\n", err)
			return true, 1
		}
		return true, 0
	case "--uninstall":
		if err := installer.Uninstall(); err != nil {
			slog.Error("uninstall failed", "err", err)
			fmt.Fprintf(os.Stderr, "[uninstall] %v\n", err)
			return true, 1
		}
		return true, 0
	case "--status":
		ui.PrintStatus()
		return true, 0
	case "--benchmark":
		if err := ui.RunBenchmarkOnce(); err != nil {
			slog.Error("benchmark failed", "err", err)
			fmt.Fprintf(os.Stderr, "[benchmark] %v\n", err)
			return true, 1
		}
		return true, 0
	}
	return false, 0
}
