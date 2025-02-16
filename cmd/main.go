package main

import (
	"MerchShop/internal/adapters/db"
	"MerchShop/internal/adapters/router"
	"MerchShop/internal/application/core/api"
	"MerchShop/internal/application/core/tokens"
	"MerchShop/internal/config"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	tokenHandler := tokens.NewTokenHandler([]byte(config.GetSecretKey()))
	dbAdapter, err := db.NewDBAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatal(err)
	}
	app := api.NewApplication(dbAdapter, tokenHandler)
	port := config.GetApplicationPort()
	rtr := router.NewRouter(app, port)
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return rtr.Start()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return rtr.Stop(context.Background())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
	log.Println("closing database connection...")
	dbAdapter.Close()
	log.Println("finished")
}
