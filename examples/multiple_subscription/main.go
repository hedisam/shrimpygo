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
	client, err := shrimpygo.NewClient(shrimpyCfg)
	if err != nil {
		log.Fatal(err)
	}

	// listen and receive data for 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ws, err := client.Websocket(ctx, 10)
	if err != nil {
		log.Fatal(err)
	}

	// a list of the subscriptions that we're interested in.
	// ordersSubs: all limit orders across all exchanges that have been completed by Shrimpy in real-time
	ordersSubs := shrimpygo.OrdersSubs()
	// tradesSubs: the executed trades for the given pair on the specified exchange
	tradesSubs := shrimpygo.TradesSubs("binance", "btc-usdt")
	// bboSubs: the BBO channel (Best-Bid-Offer) provides the best bid and ask prices for the given pair
	bboSubs := shrimpygo.BBOSubs("coinbasepro", "btc-usd")
	// orderBookSubs & obBinanceSubs: the full order book for the given pair and exchange
	orderBookSubs := shrimpygo.OrderBookSubs("coinbasepro", "btc-usd")
	obBinanceSubs := shrimpygo.Subscription{
		Type:     "subscribe",
		Exchange: "binance",
		Pair:     "btc-usdt",
		Channel:  "orderbook",
	}

	ws.Subscribe(bboSubs, orderBookSubs, ordersSubs, tradesSubs, obBinanceSubs)

	var bboUnsubscribed bool
	var counter int

	for msg := range ws.Stream() {
		fmt.Println("=========================================")
		switch message := msg.(type) {
		case *shrimpygo.OrderBookInfo: // it could be from bbo or order-book channel
			if message.Snapshot {
				// too much data to be printed.
				continue
			}
			if message.Channel == shrimpygo.ChannelBBO {
				fmt.Printf("BBO from %s: %v\n", message.Exchange, message)
				if !bboUnsubscribed && counter > 5 {
					// ok, we're done with this, no more bbo data.
					fmt.Printf("\n\n//////////////// Unsubscribing from bbo channel ////////////////\n\n")
					ws.Unsubscribe(bboSubs)
					bboUnsubscribed = true
				}
				counter++
				continue
			}

			fmt.Println("OrderBook:", message)

		case *shrimpygo.Trades:
			fmt.Println("Trades:", message)
		case *shrimpygo.Orders:
			fmt.Println("Orders:", message)
		case error:
			log.Println("error from shrimpy:", message)
			return
		default:
			fmt.Println("unwanted data returned by the stream:", message)
		}
	}
}
