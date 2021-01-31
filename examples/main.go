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

	shrimpyCfg := shrimpygo.Config{PublicKey: cfg.APIKey, PrivateKey: cfg.SecretKey}

	client, err := shrimpygo.NewShrimpyClient(shrimpyCfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws, err := client.Websocket(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ws.Subscribe(
		shrimpygo.BBOSubs("coinbasepro", "btc-usd"),
		shrimpygo.OrdersSubs(),
		)

	for msg := range ws.Stream() {
		fmt.Println("=========================================")
		switch message := msg.(type) {
		case *shrimpygo.OrderBook: // it could be from channel bbo or orderbook
			if message.Snapshot { // snapshot contains a lot of data, fills the entire screen.
				continue
			}
			if message.Channel == shrimpygo.ChannelBBO {
				fmt.Println("bbo:", message)
				continue
			}
			fmt.Println("orderbook:", message)
		case *shrimpygo.Trades:
			fmt.Println(message)
		case *shrimpygo.Orders:
			fmt.Println(message)
		case error:
			log.Println("error from shrimpy: ", message)
			return
		default:
			fmt.Println("unwanted data returned by the stream:", message)
		}
	}
}
