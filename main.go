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
	"msg-processor/internal/broker"
	"msg-processor/internal/repository"
	"msg-processor/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := app.NewConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgxCfg, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	//pgxCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	//	pgxuuid.Register(conn.TypeMap())
	//	return nil
	//}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		log.Fatal(err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	producer := broker.NewProducer(cfg.KafkaAddrs)
	s := service.NewService(repo, producer)
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

	go func() {
		timer := time.NewTimer(0)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("send to Kafka:", ctx.Err())
				return
			case <-timer.C:
				err = s.SendMsgToKafka(ctx)
				if err != nil {
					log.Println("send to Kafka:", err)
				}

				timer.Reset(time.Second)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-c
	cancel()

	downCtx, downCancel := context.WithTimeout(context.Background(), time.Second)
	defer downCancel()

	err = srv.Shutdown(downCtx)
	if err != nil {
		log.Panic(err)
	}
}
