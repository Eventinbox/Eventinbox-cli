package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen [port]",
	Short: "Forward live webhook events to localhost",
	// Hidden until the tunnel backend exists — today this only stands up a
	// local reverse proxy and never connects to EventInbox, so it's kept out
	// of --help to avoid implying a working live-forwarding workflow.
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	Example: `  eventinbox listen 3000
  eventinbox listen 3000 --path /webhooks/payments`,
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			path = "/"
		}

		target, err := url.Parse(fmt.Sprintf("http://localhost:%s", port))
		if err != nil {
			log.Fatalf("invalid port: %v", err)
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("→ %s %s\n", r.Method, r.URL.Path)
			r.URL.Path = path
			proxy.ServeHTTP(w, r)
		})

		listenAddr := ":8090"
		fmt.Printf("EventInbox tunnel listening on %s → localhost:%s%s\n", listenAddr, port, path)
		fmt.Println("Waiting for events... (Ctrl+C to stop)")

		srv := &http.Server{Addr: listenAddr, Handler: mux}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("tunnel error: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("\nTunnel stopped.")
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.Flags().String("path", "/", "Local path to forward events to")
}
