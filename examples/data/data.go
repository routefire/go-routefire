package main

import (
	"flag"
	"fmt"
	"github.com/routefire/go-routefire"
	"os"
)

func printUsage(){
	fmt.Printf("Usage: ./data -uid username -pass p@ssw0rd\n")
	os.Exit(1)
}


func main(){

	// Set up `flag` to accept a username and password via command-line arguments
	uid := flag.String("uid", "", "username")
	password := flag.String("pass", "", "password")
	flag.Parse()
	if uid == nil || password == nil {
		printUsage()
		return
	} else if len(*uid) == 0 || len(*password) == 0 {
		printUsage()
		return
	}

	// Create a new Routefire client
	client, err := routefire.New(*uid, *password)

	if err != nil {
		panic(err)
	}

	// For Routefire Core users:
	ob2, err := client.GetConsolidatedOrderBook(*uid, routefire.Btc, routefire.Usd)

	if err != nil {
		panic(err)
	} else {
		printStandardOrderBook(ob2)
	}

	// For DMA-only users:
	ob, err := client.GetConsolidatedOrderBookDMA(*uid, routefire.Btc, routefire.Usd)

	if err != nil {
		panic(err)
	} else if len(ob.Errors) == 0 {
		printDmaOrderBook(&ob.Data)
	}

	// To get balances...
	bals, err := client.GetBalances(*uid, routefire.Btc)
	if err != nil {
		panic(err)
	} else  {
		for k, v := range bals {
			k2 := padStringForPrinting(k, 10) // Pretty-print
			fmt.Printf("%s: \t%s\n", k2, v)
		}
	}
}

//
//  UTILITY FUNCTIONS
//

func padStringForPrinting(s string, n int) string {
	m := n - len(s)
	s2 := s + ""
	for i := 0; i < m; i++ {
		if i % 2 == 0 {
			s2 = " " + s2
		} else {
			s2 = s2 + " "
		}
	}
	return s2
}

// Pretty-print (DMA) order book
func printDmaOrderBook(ob *routefire.DmaOrderBook) {
	fmt.Printf("---------------- OFFERS ------------------\n")
	for i, x := range ob.Offers {
		fmt.Printf("ASK\t%d\t%s @ %s\n", i, x.Amount, x.Price)
	}
	fmt.Printf("----------------- BIDS -------------------\n")
	for i, x := range ob.Bids {
		fmt.Printf("BID\t%d\t%s @ %s\n", i, x.Amount, x.Price)
	}
	fmt.Printf("------------------------------------------\n\n\n")
}

// Pretty-print Routefire Core order book
func printStandardOrderBook(ob *routefire.OrderBookResponse) {
	fmt.Printf("---------------- OFFERS ------------------\n")
	for i, x := range ob.Offers {
		fmt.Printf("ASK\t%d\t%s @ %s\n", i, x.Quantity, x.Price)
	}
	fmt.Printf("----------------- BIDS -------------------\n")
	for i, x := range ob.Bids {
		fmt.Printf("BID\t%d\t%s @ %s\n", i, x.Quantity, x.Price)
	}
	fmt.Printf("------------------------------------------\n\n\n")
}

