package main

import (
	"context"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"iryzzh/practicum-gophermart/cmd/gophermart/config"
	"iryzzh/practicum-gophermart/internal/app/accrual"
	"iryzzh/practicum-gophermart/internal/app/handler"
	"iryzzh/practicum-gophermart/internal/app/server"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"iryzzh/practicum-gophermart/internal/app/store/pgstore"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	var s store.Store

	switch {
	case cfg.DatabaseURI != "":
		if s, err = pgstore.New(cfg.DatabaseURI); err != nil {
			log.Fatal(err)
		}
	default:
		s = memstore.New()
	}

	defer s.Close()

	h := handler.New(s, []byte(cfg.SessionKey))
	srv := server.New(cfg.RunAddress, h)

	ac := accrual.New(s, cfg.AccrualSystemAddress, time.Second)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ac.Serve(ctx)
	})

	g.Go(func() error {
		return srv.Serve(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
