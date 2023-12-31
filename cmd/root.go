package cmd

import (
	"log"
	"net/http"
	"strconv"

	"github.com/geekxflood/workspaceone-exporter/internal/httpclient"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	port        string
	insecure    bool
	defaultLgid int
)

var rootCmd = &cobra.Command{
	Use:   "workspaceone-exporter",
	Short: "A Prometheus exporter for WorkspaceOne UEM",
	Long:  `workspaceone-exporter is a tool that exposes metrics from WorkspaceOne UEM to Prometheus. It provides an HTTP server for Prometheus to scrape.`,
	Run:   runServer,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error during command execution: %v", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "9100", "Port to run the HTTP server on")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "i", false, "Ignore TLS for self-signed certificates")
	rootCmd.PersistentFlags().IntVarP(&defaultLgid, "default-lgid", "d", 0, "Default LGID value to use if not provided in query")
}

func runServer(cmd *cobra.Command, args []string) {
	client := httpclient.New(insecure)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(w, r, client)
	})

	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request, client *http.Client) {
	queryParams := r.URL.Query()
	lgidParam := queryParams.Get("lgid")

	lgid := defaultLgid
	if lgidParam != "" {
		var err error
		lgid, err = strconv.Atoi(lgidParam)
		if err != nil {
			log.Printf("Invalid lgid query parameter: %v", err)
			http.Error(w, "Invalid lgid parameter", http.StatusBadRequest)
			return
		}
	}

	// Example: Update your metrics based on the lgid
	// This is where you would call your internal API functions using the client
	// and update your Prometheus metrics

	// Serve the metrics
	promhttp.Handler().ServeHTTP(w, r)
}
