// Package cmd implements the jellyfin-cli commands.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "jellyfin",
	Short: "CLI for managing Jellyfin media server",
	Long: `CLI for managing Jellyfin media server.

Quick examples:
  jellyfin login -u user -p pass                 Authenticate
  jellyfin libraries                              List libraries
  jellyfin movies                                 List all movies
  jellyfin search "kung fu"                       Search for media
  jellyfin info <id>                              Show item details
  jellyfin update <id> --name "Title" --year 2020 Update metadata
  jellyfin identify <id> --name "Title" --imdb tt1234567  Identify/match item
  jellyfin refresh <id>                           Re-fetch metadata from providers
  jellyfin scan                                   Scan all libraries for changes
  jellyfin sessions                               Show active playback sessions`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output raw JSON responses")
	rootCmd.SetHelpFunc(customHelp)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func customHelp(cmd *cobra.Command, _ []string) {
	if cmd == rootCmd {
		printTree()
		return
	}
	printLeafHelp(cmd)
}

func printTree() {
	fmt.Println("jellyfin-cli - Jellyfin media server CLI")
	fmt.Println("")
	fmt.Println("Global: --json (raw JSON output)")
	fmt.Println("")
	fmt.Println("Commands:")

	for _, cmd := range rootCmd.Commands() {
		if cmd.Hidden || cmd.Name() == "help" || cmd.Name() == "completion" {
			continue
		}
		aliases := ""
		if len(cmd.Aliases) > 0 {
			aliases = " (" + strings.Join(cmd.Aliases, ", ") + ")"
		}
		fmt.Printf("  %-12s %s%s\n", cmd.Name(), cmd.Short, aliases)
	}

	fmt.Println("")
	fmt.Println("Run 'jellyfin <command> --help' for full details.")
}

func printLeafHelp(cmd *cobra.Command) {
	fmt.Printf("%s\n", cmd.UseLine())
	fmt.Println(cmd.Short)

	if cmd.HasLocalFlags() {
		fmt.Println("")
		fmt.Println("Flags:")
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			shorthand := ""
			if f.Shorthand != "" {
				shorthand = "-" + f.Shorthand + ", "
			}
			def := ""
			if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
				def = " (default: " + f.DefValue + ")"
			}
			fmt.Printf("  %s--%s %s%s\n", shorthand, f.Name, f.Usage, def)
		})
	}
}
