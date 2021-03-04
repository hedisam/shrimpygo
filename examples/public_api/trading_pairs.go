package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"log"
)

func TradingPairs(cli *shrimpygo.Client, freeCall bool) {
	pairs, err := cli.TradingPairs(context.Background(), "coinbasepro", freeCall)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s's pairs: %v\n", "coinbasepro", pairs)
}
