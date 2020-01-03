package routefire

import (
	"fmt"
	"log"
	"testing"
)

const (
	UnitTestOrderId = "1f8a9516-4fc8-467c-840a-901f04e3a7a3"
	uid             = "johndoe@gmail.com"
	password        = "password"
)

var apiClient, _ = New(uid, password)

func TestRouteFireAPI_GetConsolidatedOrderBook(t *testing.T) {
	resp, err := apiClient.GetConsolidatedOrderBook(uid, "btc", "eur")

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_GetOBStats(t *testing.T) {
	resp, err := apiClient.GetOrderBookStats(uid, "btc", "eur", "5.5")

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_GetBalances(t *testing.T) {
	resp, err := apiClient.GetBalances(uid, "btc")

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

	resp, err := apiClient.SubmitOrder(uid, "btc", "usd", "0.003", "", "rfxw", params)
	if err != nil {
		t.Errorf("SubmitOrder should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_SubmitOrderLimit(t *testing.T) {
	params := map[string]string{
		"target_seconds": "100",
		"backfill":       "1.0",
		"aggression":     "0.0",
		"iwould":         "8000.0", // The limit price
	}

	// NOTE: buying USD, selling BTC = selling Bitcoin for U.S. dollars
	resp, err := apiClient.SubmitOrder(uid, "usd", "btc", "0.003", "", "rfxw", params)
	if err != nil {
		t.Errorf("SubmitOrder should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_GetOrderStatus(t *testing.T) {
	resp, err := apiClient.GetOrderStatus(uid, UnitTestOrderId)

	if err != nil {
		t.Errorf("GetOrderStatus should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_CancelOrder(t *testing.T) {
	resp, err := apiClient.CancelOrder(uid, "9f2b14dc-1f67-4d0c-9270-0dfe23cc36b7")

	if err != nil {
		t.Errorf("CancelOrder should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestDmaAPI_GetConsolidatedOrderBook(t *testing.T) {
	resp, err := apiClient.GetConsolidatedOrderBookDMA(uid, Btc, Usd)

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestDmaAPI_Balances(t *testing.T) {
	rig, err := apiClient.BalanceDMA(uid, "GEMINI", Btc)

	if err != nil {
		t.Errorf("BalanceDMA should not return error, got %s\n", err)
	} else {
		log.Printf("Balance: %s %s\n", rig.Amount, rig.Asset)
	}
}
