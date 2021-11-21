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

	"github.com/qiniu/qmgo"
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
	dbURI := flag.String("db-uri", "mongodb://mongo:27017", "MongoDB URI")
	dbName := flag.String("db-name", "alpha", "MongoDB DB name")
	flag.Parse()

	client, err := qmgo.Open(ctx, &qmgo.Config{
		Uri:      *dbURI,
		Database: *dbName,
		Coll:     "attempts",
	})
	if err != nil {
		return err
	}

	srv := alpha.NewMongoServer(client)
	handler, err := srv.Handler()
	if err != nil {
		return err
	}

	httpSrv := &http.Server{
		Addr:    *addr,
		Handler: handler,
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
