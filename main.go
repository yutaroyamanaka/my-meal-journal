package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/go-kit/log"

	"github.com/spf13/cobra"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	var (
		port int
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run very simple HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := &http.Server{
				Addr: fmt.Sprintf(":%d", port),
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("hello world"))
				}),
			}
			logger.Log("msg", "http server starts runnning", "port", port)
			return s.ListenAndServe()
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 80, "port number that http server runs on")
	if err := cmd.Execute(); err != nil {
		logger.Log("error", err)
	}
}
