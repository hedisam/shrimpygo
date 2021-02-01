package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"github.com/hedisam/shrimpygo/examples/appconfig"
	"log"
	"time"
)

func main() {
	cfg, err := appconfig.Read("examples/appconfig/config.json")
	if err != nil {
		log.Fatal(err)
	}

	shrimpyCfg := shrimpygo.Config{PublicKey: cfg.APIKey, PrivateKey: cfg.SecretKey}
	client, err := shrimpygo.NewShrimpyClient(shrimpyCfg)
	if err != nil {
		log.Fatal(err)
	}

	// receive data for 5 seconds, unless there's an error.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()

	// you can specify the stream's throughput (the data channel's capacity)
	// by setting throughput as the the second parameter.
	ws, err := client.Websocket(ctx, 0)
	if err != nil {
		log.Fatal(err)
	}

	// subscribing to the order-book channel.
	// check the examples to see how you can subscribe to other channels.
	ws.Subscribe(
		shrimpygo.OrderBookSubs("coinbasepro", "btc-usd"),
	)

	// reading the stream which can push any type of data from the supported
	// ws channels (bbo, trades, etc); depends on how many different channels
	// you have subscribed on the same ws connection.
	for iData := range ws.Stream() {
		fmt.Println("============================")
		switch data := iData.(type) {
		case *shrimpygo.OrderBook:
			if data.Snapshot {
				// too much data to be printed.
				continue
			}
			fmt.Printf("OrderBook: %v\n", data)
		case error:
			fmt.Printf("shrimpy error: %+v\n", data)
			return
		default:
			fmt.Printf("unknown/unwanted data received from the stream: %v\n", data)
		}
	}
}
