package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:     "send [event-type]",
	Short:   "Send a test event to an endpoint",
	Args:    cobra.ExactArgs(1),
	Example: `  eventinbox send payment.created --tenant acme --endpoint payments --payload '{"amount":5400}'`,
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Root().PersistentFlags().GetString("api-url")
		eventType := args[0]
		tenant, _ := cmd.Flags().GetString("tenant")
		endpoint, _ := cmd.Flags().GetString("endpoint")
		payload, _ := cmd.Flags().GetString("payload")

		if tenant == "" || endpoint == "" {
			fmt.Println("error: --tenant and --endpoint are required")
			os.Exit(1)
		}
		if payload == "" {
			payload = `{}`
		}

		url := fmt.Sprintf("%s/in/%s/%s", apiURL, tenant, endpoint)
		req, _ := http.NewRequest("POST", url, bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-ei-event-type", eventType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusAccepted {
			fmt.Printf("✓ Event sent (%s)\n%s\n", eventType, string(body))
		} else {
			fmt.Printf("error: server returned %d\n%s\n", resp.StatusCode, string(body))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().String("tenant", "", "Tenant name")
	sendCmd.Flags().String("endpoint", "", "Endpoint name")
	sendCmd.Flags().String("payload", "{}", "JSON payload")
}
