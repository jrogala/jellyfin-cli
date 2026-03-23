// Package cmdutil provides shared helpers for CLI commands.
package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jrogala/jellyfin-cli/client"
	"github.com/jrogala/jellyfin-cli/config"
	"github.com/spf13/cobra"
)

// NewClient creates an authenticated Jellyfin client from config.
func NewClient() (*client.Client, error) {
	session, err := config.LoadSession()
	if err != nil {
		return nil, err
	}
	return client.New(config.URL(), session.Token, session.UserID), nil
}

// PrintJSON encodes v as indented JSON to stdout.
func PrintJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// IsJSON returns true if the --json persistent flag is set on the command's root.
func IsJSON(cmd *cobra.Command) bool {
	v, _ := cmd.Root().PersistentFlags().GetBool("json")
	return v
}

// Render outputs data as JSON if --json is set, otherwise calls tableFunc.
func Render(cmd *cobra.Command, data any, tableFunc func()) {
	if IsJSON(cmd) {
		_ = PrintJSON(data)
		return
	}
	tableFunc()
}

// NewTabWriter returns a tabwriter configured for CLI output.
func NewTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

// ExitErr prints an error to stderr and exits with code 1.
func ExitErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
