package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"msg-processor/internal/api"
	"msg-processor/internal/app"
	"msg-processor/internal/repository"
	"msg-processor/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := app.NewConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	s := service.NewService(repo)
	h := api.NewHandler(s)

	r := http.NewServeMux()
	r.HandleFunc("POST /messages", h.CreateMsg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		err = srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-c

	downCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = srv.Shutdown(downCtx)
	if err != nil {
		log.Panic(err)
	}
}
