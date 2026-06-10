package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var replayCmd = &cobra.Command{
	Use:     "replay [delivery-id]",
	Short:   "Replay a delivery by ID",
	Args:    cobra.ExactArgs(1),
	Example: `  eventinbox replay b4901ce3-d019-4c5c-96a2-d7990b045b7b`,
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Root().PersistentFlags().GetString("api-url")
		apiKey, _ := cmd.Root().PersistentFlags().GetString("api-key")
		workspaceID, _ := cmd.Root().PersistentFlags().GetString("workspace")
		if apiKey == "" {
			apiKey = os.Getenv("EI_API_KEY")
		}
		if workspaceID == "" {
			workspaceID = os.Getenv("EI_WORKSPACE_ID")
		}
		deliveryID := args[0]

		if apiKey == "" || workspaceID == "" {
			fmt.Println("error: --api-key and --workspace are required")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/api/v1/deliveries/%s/replay", apiURL, deliveryID)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-workspace-id", workspaceID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusAccepted {
			fmt.Printf("✓ Delivery %s queued for replay\n", deliveryID)
		} else {
			fmt.Printf("error: server returned %d\n", resp.StatusCode)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)
}
