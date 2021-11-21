package main

import (
	"alpha"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	addr := flag.String("addr", ":8080", "HTTP server listen address")
	flag.Parse()

	srv := alpha.NewServer()
	mux := http.NewServeMux()
	h, err := srv.HandleIndex()
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	mux.HandleFunc("/", h)
	mux.HandleFunc("/increment", srv.HandleIncrement())

	httpSrv := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		shutdownCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		errCh <- httpSrv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("listening on %s\n", *addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
