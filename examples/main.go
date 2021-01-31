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

	ws.Subscribe(shrimpygo.ChanBBO, "coinbasepro", "btc-usd")

	for msg := range ws.Stream() {
		fmt.Println("=========================================")
		switch message := msg.(type) {
		case *shrimpygo.OrderBook:
			if message.Snapshot { // snapshot contains a lot of data, fills the entire screen.
				continue
			}
			fmt.Println(message)
		case *shrimpygo.Trades:
			fmt.Println(message)
		case error:
			log.Println("error from shrimpy: ", message)
			return
		default:
			fmt.Println("unwanted data returned by the stream:", message)
		}
	}
}
