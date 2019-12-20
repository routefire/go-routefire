package main

import (
	"flag"
	"fmt"
	"github.com/routefire/go-routefire"
	"os"
	"time"
)


func printUsage() {
	fmt.Printf("Usage: ./momtrader -uid username -pass p@ssw0rd\n")
	os.Exit(1)
}
func main(){

	// Set up `flag` to accept a username and password via command-line arguments
	uid := flag.String("uid", "", "username")
	password := flag.String("pass", "", "password")
	flag.Parse()

	// Check all the inputs are valid.
	if uid == nil || password == nil{
		printUsage()
		return
	} else if len(*uid) == 0 || len(*password) == 0  {
		printUsage()
		return
	}
	// Create a new Routefire client.
	client, err := routefire.New(*uid, *password)
	if err != nil {
		panic(err)
	}

	assets := []string{"btc","eth","zrx"}
	trader := NewMomentumTrader(*uid, client, assets, "usd", 40.0, 1.0) // Trade with $100
	trader.RunLoop(10 * time.Second)
}


