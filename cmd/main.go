package main

import (
	"MerchShop/internal/adapters/db"
	"MerchShop/internal/adapters/router"
	"MerchShop/internal/application/core/api"
	"MerchShop/internal/application/core/tokens"
	"MerchShop/internal/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	if config.GetEnv() == "development" {
		log.SetLevel(log.DebugLevel)
	}
	tokenHandler := tokens.NewTokenHandler([]byte(config.GetSecretKey()))
	dbAdapter, err := db.NewDBAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatal(err)
	}
	app := api.NewApplication(dbAdapter, tokenHandler)
	rtr := router.NewRouter(app, config.GetApplicationPort())
	err = rtr.Start()
	if err != nil {
		log.Fatal(err)
	}
}
