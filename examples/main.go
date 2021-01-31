package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"github.com/hedisam/shrimpygo/examples/config"
	"log"
)

func main() {
	cfg, err := config.Read("examples/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	client, err := shrimpygo.NewShrimpyClient(cfg.APIKey, cfg.SecretKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws, err := client.Websocket(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ws.Subscribe("orderbook", "coinbasepro", "btc-usd")

	for msg := range ws.Stream() {
		switch message := msg.(type) {
		case shrimpygo.PriceData:
			if message.Snapshot {
				continue
			}
			fmt.Println("=========================================")
			fmt.Println(message)
		case error:
			log.Println("error from shrimpy: ", message)
			return
		default:
			log.Println("unknown data type")
		}
	}
}
