// package main sets up the http server.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/spf13/cobra"
	"github.com/yutaroyamanaka/my-meal-journal/internal/clock"
	"github.com/yutaroyamanaka/my-meal-journal/internal/handler"
	"github.com/yutaroyamanaka/my-meal-journal/internal/service"
	"github.com/yutaroyamanaka/my-meal-journal/internal/store"
)

var port int

func newLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return logger
}

func run(ctx context.Context, db *sql.DB, l net.Listener, logger log.Logger) error {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		level.Info(logger).Log("msg", "health check is conducted")
		w.Write([]byte("OK\n"))
	}))

	r, err := store.NewRepository(&clock.RealTimeClocker{}, db)
	if err != nil {
		return err
	}
	addsvc, err := service.NewAddService(r, logger)
	if err != nil {
		return err
	}
	addh, err := handler.NewAddHandler(addsvc, logger)
	if err != nil {
		return err
	}
	mux.Handle("/add", addh)

	srv := &http.Server{
		Handler: mux,
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
		Short: "Run an API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := store.Config{}
			if err := env.Parse(&c); err != nil {
				return err
			}
			db, cleanup, err := store.Open(&c)
			if err != nil {
				return err
			}
			defer cleanup()

			l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				level.Error(logger).Log("msg", fmt.Sprintf("failed to listen on address: %d", port), "error", err)
				return err
			}
			return run(ctx, db, l, logger)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 80, "port number that http server runs on")
	return cmd
}

func main() {
	logger := newLogger()
	cmd := app(context.Background(), logger)
	if err := cmd.Execute(); err != nil {
		level.Error(logger).Log("msg", "failed to execute a command", "error", err)
		os.Exit(1)
	}
}
