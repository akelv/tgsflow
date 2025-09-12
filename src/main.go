package main

import (
	"flag"
	"fmt"
	"os"

	cmdpkg "github.com/kelvin/tgsflow/src/cmd"
	"github.com/kelvin/tgsflow/src/util/logx"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printVersion() {
	fmt.Printf("tgs %s (commit %s, built %s)\n", version, commit, date)
}

func main() {
	// Global flags
	var (
		showVersion bool
		jsonLogs    bool
	)

	fs := flag.NewFlagSet("tgs", flag.ContinueOnError)
	fs.BoolVar(&showVersion, "version", false, "Print version and exit")
	fs.BoolVar(&jsonLogs, "json", false, "Emit JSONL logs to stderr")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	if jsonLogs {
		logx.SetJSON(true)
	}

	if showVersion {
		printVersion()
		return
	}

	args := fs.Args()
	if len(args) == 0 {
		cmdpkg.CmdHelp(nil)
		return
	}

	sub := args[0]
	subArgs := args[1:]

	switch sub {
	case "help", "--help", "-h":
		cmdpkg.CmdHelp(subArgs)
	case "init":
		os.Exit(cmdpkg.CmdInit(subArgs))
	case "context":
		os.Exit(cmdpkg.CmdContext(subArgs))
	case "specify":
		os.Exit(cmdpkg.CmdSpecify(subArgs))
	case "plan":
		os.Exit(cmdpkg.CmdPlan(subArgs))
	case "tasks":
		os.Exit(cmdpkg.CmdTasks(subArgs))
	case "approve":
		os.Exit(cmdpkg.CmdApprove(subArgs))
	case "verify":
		os.Exit(cmdpkg.CmdVerify(subArgs))
	case "brief":
		os.Exit(cmdpkg.CmdBrief(subArgs))
	case "version":
		printVersion()
	case "implement", "drift-check", "docs", "open-pr", "watch":
		// Stubs for now
		logx.Infof("Command '%s' is not implemented yet in this milestone.", sub)
		os.Exit(2)
	default:
		logx.Errorf("Unknown subcommand: %s", sub)
		cmdpkg.CmdHelp(nil)
		os.Exit(2)
	}
}
