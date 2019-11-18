package routefire

import (
	"fmt"
	"testing"
)

const (
	UnitTestOrderId = "1f8a9516-4fc8-467c-840a-901f04e3a7a3"
	uid = "johndoe@gmail.com"
	password = "password"
)

var routeFireApi, _ = New(uid, password)

func TestRouteFireAPI_GetConsolidatedOrderBook(t *testing.T) {
	resp, err := routeFireApi.GetConsolidatedOrderBook(uid, "btc", "eur")

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_GetOBStats(t *testing.T) {
	resp, err := routeFireApi.GetOrderBookStats(uid, "btc", "eur", "5.5")

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_GetBalances(t *testing.T) {
	resp, err := routeFireApi.GetBalances(uid, "btc")

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", resp)
}

func TestRouteFireAPI_SubmitOrder(t *testing.T) {
	params := map[string]string{
		"target_seconds": "100",
		"backfill":       "1.0",
		"aggression":     "0.0",
	}

	resp, err := routeFireApi.SubmitOrder(uid, "btc", "usd", "0.003", "10000", "rfxw", params)
	if err != nil {
		t.Errorf("SubmitOrder should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}
func TestRouteFireAPI_GetOrderStatus(t *testing.T) {
	resp, err := routeFireApi.GetOrderStatus(uid, UnitTestOrderId)

	if err != nil {
		t.Errorf("GetOrderStatus should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_CancelOrder(t *testing.T) {
	resp, err := routeFireApi.CancelOrder(uid, "9f2b14dc-1f67-4d0c-9270-0dfe23cc36b7")

	if err != nil {
		t.Errorf("CancelOrder should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestDmaAPI_GetConsolidatedOrderBook(t *testing.T) {
	resp, err := routeFireApi.GetConsolidatedOrderBookDMA(uid, Btc, Usd)

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}


