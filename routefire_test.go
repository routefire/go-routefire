package routefireApi

import (
	"fmt"
	"testing"
)

const uid = "johndoe@gmail.com"
const password = "password"

var routeFireApi = New(uid, password)

func TestRouteFireAPI_GetConsolidatedOrderBook(t *testing.T) {
	resp, err := routeFireApi.GetConsolidatedOrderBook(uid, "btc", "eur", "10.34")

	if err != nil {
		t.Errorf("GetConsolidatedOrderBook should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_SubmitOrder(t *testing.T) {
	params := map[string]string{
		"target_seconds": "300",
		"backfill":       "1.0",
		"aggregation":    "0.0",
	}

	resp, err := routeFireApi.SubmitOrder(uid, "btc", "eur", "9.34", "10000", "vwap", params)
	if err != nil {
		t.Errorf("SubmitOrder should not return error, got %s\n", err)
	}

	fmt.Printf("%+v\n", *resp)
}

func TestRouteFireAPI_GetOrderStatus(t *testing.T) {
	resp, err := routeFireApi.GetOrderStatus(uid, "9f2b14dc-1f67-4d0c-9270-0dfe23cc36b7")

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
