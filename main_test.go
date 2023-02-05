package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"

	log "github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
)

func TestRun_health_check(t *testing.T) {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen on local address: %#v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error, 1)
	go func() {
		ch <- run(ctx, l, logger)
	}()
	url := fmt.Sprintf("http://%s/health", l.Addr().String())
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}
	want := "OK\n"
	if diff := cmp.Diff(string(got), want); diff != "" {
		t.Errorf("Got an unexcpeted response: %s", diff)
	}
	cancel()
	if err := <-ch; err != nil {
		t.Errorf("run func returned an error: %#v", err)
	}
}
