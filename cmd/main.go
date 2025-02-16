package main

import (
	"MerchShop/internal/adapters/db"
	"MerchShop/internal/adapters/router"
	"MerchShop/internal/application/core/api"
	"MerchShop/internal/application/core/tokens"
	"MerchShop/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	if config.GetEnv() == "development" {
		log.SetLevel(log.DebugLevel)
	}
	tokenHandler := tokens.NewTokenHandler([]byte(config.GetSecretKey()))
	dbAdapter, err := db.NewDBAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatal(err)
	}
	app := api.NewApplication(dbAdapter, tokenHandler)
	port := config.GetApplicationPort()
	rtr := router.NewRouter(app, port)
	go func() {
		log.Printf("Server is starting on :%s...", port)
		err := rtr.Start()
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	<-stop
	log.Println("closing database connection...")
	dbAdapter.Close()
	log.Println("finished")
}
