package main

import (
	"flag"
	"fmt"
	"github.com/routefire/go-routefire"
	"os"
	"strconv"
	"time"
)

func printUsage() {
	fmt.Printf("Usage: ./dma -uid username -pass p@ssw0rd -quantity 1.234\n")
	os.Exit(1)
}

func main() {

	// Set up `flag` to accept a username and password via command-line arguments
	uid := flag.String("uid", "", "username")
	password := flag.String("pass", "", "password")
	quantity := flag.String("quantity", "", "quantity")
	flag.Parse()

	// Check all the inputs are valid.
	if uid == nil || password == nil || quantity == nil {
		printUsage()
		return
	} else if len(*uid) == 0 || len(*password) == 0 || len(*quantity) == 0 {
		printUsage()
		return
	}
	checkIsFloat(*quantity)

	// Create a new Routefire client.
	client, err := routefire.New(*uid, *password)

	if err != nil {
		panic(err)
	}

	// Look up the current order book.
	obData, err := client.GetConsolidatedOrderBookDMA(*uid, routefire.Btc, routefire.Usd)

	if err != nil {
		panic(err)
	}

	// Calculate the price to submit. We'll use 1 penny less than the current mid price.
	ourPx := calcMidPrice(obData, -0.01)

	// Submit an order to buy at the best-offered venue at 0.01 less than the mid price from
	// the _consolidated_ order book.
	bestOffer := obData.Data.Offers[0] // They are sorted from lowest to highest.
	bestVenue := bestOffer.Venue

	fmt.Printf("Submitting to venue %s at price level %s\n", bestVenue, ourPx)

	// And submit the order.
	customParams := map[string]string{} // No custom parameters
	orderConfirm, err := client.SubmitOrderDMA(*uid, bestVenue, routefire.Btc, routefire.Usd, routefire.SideBuy, *quantity, ourPx, customParams)

	// Check for client error...
	if err != nil {
		panic(err)
	} else if len(orderConfirm.Errors) > 0 {
		// ... and for venue errors.
		printErrors("submit", orderConfirm.Errors)
		return
	} else {
		fmt.Printf("Successfully submitted to venue %s: order ID %s\n", orderConfirm.VenueId, orderConfirm.VenueOrderId)
	}

	// Now that we have the order ID, we can check its status. We'll do this in a loop until we
	// give up after 100 iterations (5 minutes).
	var i int
	var status *routefire.DmaOrderStatusResponse
	totalIters := 100
	for i = 0; i < totalIters; i++ {
		orderId := orderConfirm.VenueOrderId
		status, err = client.OrderStatusDMA(*uid, bestVenue, orderId)
		fmt.Printf("\n\nStatus: %+v\n\n", status)

		if err != nil {
			panic(err)
		} else if len(status.Errors) > 0 {
			printErrors("status", status.Errors)
			//return
		}

		if status.Status == routefire.StatusOpen || status.Status == routefire.StatusPartiallyFilled || len(status.Status) == 0 {
			fmt.Printf("Still working. %s filled so far.\n", status.FilledAmount)
		} else {
			fmt.Printf("Order in completed state: %s\n", status.Status)
			break
		}
		time.Sleep(3*time.Second)
	}

	// If we went all 100 iterations and the order still hasn't filled, we'll go ahead and cancel it.
	if i == totalIters {
		fmt.Printf("Canceling order at %s (ID: %s)...\n", bestVenue, orderConfirm.VenueOrderId)
		response, err := client.CancelOrderDMA(*uid, bestVenue, orderConfirm.VenueOrderId)

		if err != nil {
			panic(err)
		} else if len(response.Errors) > 0 {
			printErrors("cancel", response.Errors)
		} else {
			fmt.Printf("Cancel successful.\n")
		}
	} else {
		fmt.Printf("Order filled successfully: total %s (%s).\n", status.FilledAmount, status.Status)
	}

}

//
//  UTILITY FUNCTIONS
//

func printErrors(name string, errs []routefire.DmaError) {
	fmt.Printf("Errors on %s:\n", name)
	for i, e := range errs {
		fmt.Printf("%d: %s", i, e.Message)
	}
	os.Exit(1)
}

func calcMidPrice(obData *routefire.DmaOrderBookResponse, adjustment float64) string {
	bestOffer := obData.Data.Offers[0] // They are sorted from lowest to highest
	bestBid := obData.Data.Bids[len(obData.Data.Bids)-1]
	bestOfferFloat, _ := strconv.ParseFloat(bestOffer.Price, 64)
	bestBidFloat, _ := strconv.ParseFloat(bestBid.Price, 64)
	midPxFloat := ((bestOfferFloat + bestBidFloat) / 2.0) + adjustment
	midPxString := fmt.Sprintf("%.2f", midPxFloat)
	return midPxString
}

func checkIsFloat(s string) {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
}

