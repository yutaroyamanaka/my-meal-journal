package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/spf13/cobra"
)

var (
	port int
)

func NewLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return logger
}

func app(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Book API",
		Short: "Run a Book API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				level.Error(logger).Log("msg", fmt.Sprintf("failed to listen on address: %d", port), "error", err)
				return err
			}
			srv := &http.Server{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("Hello world\n"))
				}),
			}

			srvCh := make(chan error, 1)
			go func() {
				level.Info(logger).Log("msg", "http server starts runnning", "port", port)
				if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
					level.Error(logger).Log("msg", "failed to start the server", "error", err)
					srvCh <- err
				}
				close(srvCh)
			}()

			return <-srvCh
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 80, "port number that http server runs on")
	return cmd
}

func main() {
	logger := NewLogger()
	cmd := app(logger)
	if err := cmd.Execute(); err != nil {
		level.Error(logger).Log("msg", "failed to execute a command", "error", err)
		os.Exit(1)
	}
}
