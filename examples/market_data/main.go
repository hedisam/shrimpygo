package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"github.com/hedisam/shrimpygo/examples/appconfig"
	"log"
)

func main() {
	client := NewClient()
	RetrieveMultipleOrderBooks(client)
}

func RetrieveMultipleOrderBooks(cli *shrimpygo.Client) {
	orderBooks, err := cli.GetOrderBooks(context.Background(), "coinbasepro",
		shrimpygo.QueryParams(shrimpygo.BaseSymbol, "BTC,ETH"),
		shrimpygo.QueryParams("quoteSymbol", "USD"),
		"limit=5", // or you could just type your query without utilizing the helper function
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Market OrderBooks: \n%v\n", orderBooks)
}

func RetrieveOrderBooks(cli *shrimpygo.Client) {
	orderBooks, err := cli.GetOrderBooks(context.Background(), "coinbasepro")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Market OrderBooks: \n%v\n", orderBooks)
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