package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PraveenGongada/catalyst/internal/tui"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func main() {
	configPath := flag.String(
		"config",
		"",
		"Path to the configuration file (default: $CATALYST_CONFIG or ./catalyst.yaml)",
	)

	versionFlag := flag.Bool(
		"version",
		false,
		"Print version information and exit",
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Catalyst %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
		os.Exit(0)
	}

	if err := tui.Start(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting Catalyst: %v\n", err)
		os.Exit(1)
	}
}
