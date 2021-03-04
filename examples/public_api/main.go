package main

import (
	"github.com/hedisam/shrimpygo"
	"github.com/hedisam/shrimpygo/examples/appconfig"
	"log"
)

func main() {
	TradingPairs(NewClient(), true)
}

func NewClient() *shrimpygo.Client {
	cfg, err := appconfig.Read("examples/appconfig/config.json")
	if err != nil {
		log.Fatal(err)
	}

	shrimpyCfg := shrimpygo.Config{PublicKey: cfg.APIKey, PrivateKey: cfg.SecretKey}
	client, err := shrimpygo.NewClient(shrimpyCfg)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
