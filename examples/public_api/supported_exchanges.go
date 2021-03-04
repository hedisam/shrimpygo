package main

import (
	"context"
	"fmt"
	"log"
)

func SupportedExchanges(freeApiCall bool) {
	cli := NewClient()
	exchanges, err := cli.SupportedExchanges(context.Background(), freeApiCall)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Supported Exchanges:\n%v", exchanges)
}

