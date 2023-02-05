package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

func run(ctx context.Context, l net.Listener, logger log.Logger) error {
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello world\n"))
		}),
	}

	srvCh := make(chan error, 1)
	go func() {
		level.Info(logger).Log("msg", "http server starts runnning", "address", l.Addr().String())
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "failed to start the server", "error", err)
			srvCh <- err
		}
		close(srvCh)
	}()

	sigctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	for {
		select {
		case err := <-srvCh:
			return err
		case <-sigctx.Done():
			if err := srv.Shutdown(ctx); err != nil {
				level.Error(logger).Log("msg", "failed to shut down the server", "error", err)
				return err
			}
			return nil
		}
	}
}

func app(ctx context.Context, logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Run a Book API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				level.Error(logger).Log("msg", fmt.Sprintf("failed to listen on address: %d", port), "error", err)
				return err
			}
			return run(ctx, l, logger)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 80, "port number that http server runs on")
	return cmd
}

func main() {
	logger := NewLogger()
	cmd := app(context.Background(), logger)
	if err := cmd.Execute(); err != nil {
		level.Error(logger).Log("msg", "failed to execute a command", "error", err)
		os.Exit(1)
	}
}
