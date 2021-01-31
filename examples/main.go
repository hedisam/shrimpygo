package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"log"
)

const (
	// todo: read from a config file or from env
	apiKey    = ""
	secretKey = ""
)

func main() {
	client, err := shrimpygo.NewShrimpyClient(apiKey, secretKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws, err := client.Websocket(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ws.Subscribe("bbo", "coinbasepro", "btc-usd")

	for msg := range ws.Stream() {
		switch message := msg.(type) {
		case shrimpygo.PriceData:
			if message.Snapshot {
				continue
			}
			fmt.Println(message)
		case error:
			log.Println("error from shrimpy: ", message)
			return
		default:
			log.Println("unknown data type")
		}
	}
}
