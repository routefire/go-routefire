package main

import (
	"errors"
	"github.com/routefire/go-routefire"
	"log"
	"time"
)

const (
	DevelopmentExecutionSafety = false
)

func SubmitAndWait(uid string, client *routefire.Client, isBuy bool, asset, base, venue, quantity, price string, c chan error, waitDur time.Duration) ( *routefire.PlaceDmaOrderResponse, error) {

	if DevelopmentExecutionSafety {
		return nil, errors.New("SafetyOn")
	}

	params := map[string]string{}
	side := routefire.SideBuy
	if !isBuy {
		side = routefire.SideSell
	}
	res, err := client.SubmitOrderDMA(uid, venue, asset, base, side, quantity, price, params)
	if err != nil {
		return nil, err
	}
	oid := res.VenueOrderId
	log.Printf("For order %s %s %s/%s @ %s (%s) - got OID %s - %s\n", side, quantity, asset, base, price, venue, oid, res.VenueId)
	go func(venu, ordId string){
		for i := 0; i < 10; i++ {
			stat, err := client.OrderStatusDMA(uid, venu, ordId)
			log.Printf("Order status %s %s %s/%s @ %s (%s) (OID %s) - %s / %s\n", side, quantity, asset, base, price, venue, oid, stat.Status, stat.FilledAmount)
			if err != nil {
				c <- err
				return
			}
			if stat.Status == routefire.StatusComplete || stat.Status == routefire.StatusFilled {
				c <- nil
				return
			}
			time.Sleep(waitDur)
		}
		c <- errors.New("CouldNotComplete")
	}(venue, oid)

	return res, nil
}