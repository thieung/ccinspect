package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "ccinspect",
	Short: "Scan and inventory Claude Code configurations",
	Long:  "ccinspect scans filesystem for .claude/ directories, parses config entities (skills, hooks, agents, commands, rules, MCP, teams), and outputs a summary inventory.",
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("ccinspect v%s\n", version))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
