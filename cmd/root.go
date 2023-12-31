package cmd

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/geekxflood/workspaceone-exporter/internal/httpclient"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	port        string
	insecure    bool
	defaultLgid int
	ws1URL      string
	ws1Interval string
	tagParsing  bool
	tagFilter   string
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
	rootCmd.PersistentFlags().StringVar(&ws1URL, "ws1-url", "", "WorkspaceOne UEM base API URL endpoint")
	rootCmd.PersistentFlags().StringVar(&ws1Interval, "ws1-interval", "60", "Interval between each WS1 check to its enrolled devices in minutes")
	rootCmd.PersistentFlags().BoolVar(&tagParsing, "tag-parsing", false, "Enable or disable the tag parsing")
	rootCmd.PersistentFlags().StringVar(&tagFilter, "tag-filter", "", "String to filter Tag by it")
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

	if intervalParam := queryParams.Get("ws1-interval"); intervalParam != "" {
		ws1Interval = intervalParam
	}

	if urlParam := queryParams.Get("ws1-url"); urlParam != "" {
		if _, err := url.ParseRequestURI(urlParam); err != nil {
			log.Printf("Invalid ws1-url parameter: %v", err)
			http.Error(w, "Invalid ws1-url parameter", http.StatusBadRequest)
			return
		}
		ws1URL = urlParam
	}

	if parsingParam := queryParams.Get("tag-parsing"); parsingParam != "" {
		parsedValue, err := strconv.ParseBool(parsingParam)
		if err == nil {
			tagParsing = parsedValue
		}
	}

	if filterParam := queryParams.Get("tag-filter"); filterParam != "" {
		tagFilter = filterParam
	}

	// Here, add logic to use ws1URL, ws1Interval, tagParsing, and tagFilter
	// to fetch and update metrics

	// Serve the metrics
	promhttp.Handler().ServeHTTP(w, r)
}
