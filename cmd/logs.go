package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tail live delivery logs",
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
		tail, _ := cmd.Flags().GetBool("tail")

		if apiKey == "" || workspaceID == "" {
			fmt.Println("error: --api-key and --workspace are required")
			os.Exit(1)
		}

		fmt.Println("EventInbox delivery logs (Ctrl+C to stop)")
		fmt.Println("─────────────────────────────────────────")

		seen := map[string]bool{}

		for {
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/deliveries?limit=20", apiURL), nil)
			req.Header.Set("x-api-key", apiKey)
			req.Header.Set("x-workspace-id", workspaceID)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				time.Sleep(3 * time.Second)
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var deliveries []map[string]any
			json.Unmarshal(body, &deliveries)

			for _, d := range deliveries {
				id := fmt.Sprintf("%v", d["id"])
				if !seen[id] {
					seen[id] = true
					status := d["status"]
					icon := "●"
					if status == "delivered" {
						icon = "✓"
					} else if status == "failed" {
						icon = "✗"
					}
					fmt.Printf("%s [%s] delivery=%s event=%s\n",
						icon, status, id, d["event_id"])
				}
			}

			if !tail {
				break
			}
			time.Sleep(2 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().Bool("tail", false, "Keep polling for new deliveries")
}
