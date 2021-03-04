package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"log"
)

func ExchangeAssets(client *shrimpygo.Client, freeCall bool) {
	assets, err := client.ExchangeAssets(context.Background(), "binance", freeCall)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s's assets: %v\n", "binance", assets)
}
