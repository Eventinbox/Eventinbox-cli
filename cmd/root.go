package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eventinbox",
	Short: "EventInbox CLI — webhook delivery infrastructure",
	Long:  `Send test events, list and inspect deliveries, replay deliveries, and tail live delivery logs.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("api-url", "https://api.eventinbox.pro", "API base URL")
	rootCmd.PersistentFlags().String("api-key", "", "API key (or set EI_API_KEY env var)")
	rootCmd.PersistentFlags().String("workspace", "", "Workspace ID (or set EI_WORKSPACE_ID env var)")
}
