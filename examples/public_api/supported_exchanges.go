package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"log"
)

func SupportedExchanges(client *shrimpygo.Client,freeApiCall bool) {
	exchanges, err := client.SupportedExchanges(context.Background(), freeApiCall)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Supported Exchanges:\n%v", exchanges)
}

