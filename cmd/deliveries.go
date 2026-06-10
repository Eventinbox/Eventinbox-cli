package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var deliveriesCmd = &cobra.Command{
	Use:   "deliveries",
	Short: "List and inspect deliveries",
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
		status, _ := cmd.Flags().GetString("status")

		if apiKey == "" || workspaceID == "" {
			fmt.Println("error: --api-key and --workspace are required (or set EI_API_KEY and EI_WORKSPACE_ID)")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/api/v1/deliveries?limit=20", apiURL)
		if status != "" {
			url += "&status=" + status
		}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-workspace-id", workspaceID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var deliveries []map[string]any
		if err := json.Unmarshal(body, &deliveries); err != nil {
			fmt.Println(string(body))
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tATTEMPTS\tCREATED")
		for _, d := range deliveries {
			fmt.Fprintf(w, "%s\t%s\t%.0f\t%s\n",
				d["id"], d["status"], d["attempt_count"], d["created_at"])
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(deliveriesCmd)
	deliveriesCmd.Flags().String("status", "", "Filter by status (pending, delivered, failed)")
}
