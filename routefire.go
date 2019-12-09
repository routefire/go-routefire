// Package routefire provides access to Routefire core APIs.
package routefire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	APIURL = "https://routefire.io/api"
	ApiUrlOverrideKey = "GORF_URL"
	AdaptAPIURL = "https://routefire.io/adapt"
	APIVersion   = "v1"
	APIUserAgent = "RouteFire API client agent"
	AuthInterval = 100
)

func AdaptApiUrl() string {
	s := AdaptAPIURL
	if ev := os.Getenv(ApiUrlOverrideKey); len(ev) > 0 {
		s = strings.Replace(s, "https://routefire.io", ev, -1)
	}
	return s
}

func ApiUrl() string {
	s := APIURL
	if ev := os.Getenv(ApiUrlOverrideKey); len(ev) > 0 {
		s = strings.Replace(s, "https://routefire.io", ev, -1)
	}
	return s
}

var webHttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 10,
	},
	Timeout: time.Second * 20,
}

// Type Client provides access to a single Routefire user account.
type Client struct {
	username    string
	password    string
	accessToken string
	client      *http.Client
}

// Function New creates a new Routefire client from username/password credentials.
func New(uid, password string) (*Client, error) {
	z := &Client{uid, password, "", webHttpClient}
	if err := z.refreshToken(); err != nil {
		return nil, err
	}

	go z.refreshLoop(AuthInterval * time.Second)
	return z, nil
}

// Function SubmitOrder submits a Routefire (algorithm) order.
func (api *Client) SubmitOrder(userId string, buyAsset string, sellAsset string, quantity string, price string, algo string, algoParams map[string]string) (*SubmitOrderResponse, error) {
	var jsonData SubmitOrderResponse

	params := map[string]interface{}{
		"user_id":     userId,
		"buy_asset":   buyAsset,
		"sell_asset":  sellAsset,
		"quantity":    quantity,
		"price":       price,
		"algo":        algo,
		"algo_params": algoParams,
	}

	resp, err := api.queryPrivate("orders/submit", params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &jsonData)
	if err != nil {
		return nil, err
	}

	return &jsonData, nil
}

// Function SubmitOrderDMA submits a DMA (direct market access) order -- that is, an
// order submitted directly to a given trading venue.
func (api *Client) SubmitOrderDMA(userId, venue, asset, baseAsset string, side string, quantity, price string, orderParams map[string]string) (*PlaceDmaOrderResponse, error) {
	var jsonData PlaceDmaOrderResponse

	req := PlaceDmaOrderRequest{
		UserId:      userId,
		VenueId:     venue,
		Side:        side,
		TradedAsset: asset,
		BaseAsset:   baseAsset,
		Quantity:    quantity,
		Price:       price,
		OrderParams: orderParams,
	}

	if bs, err := json.Marshal(&req); err != nil {
		return nil, err
	} else {
		resp, err := api.queryAdaptPrivateWithBytes("orders/new", bs)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resp, &jsonData)
		if err != nil {
			return nil, err
		}

		return &jsonData, nil
	}

}

// Function OrderStatusDMA gets order status and fill amount from a given venue and order ID.
func (api *Client) OrderStatusDMA(userId, venue, venueOrdId string) (*DmaOrderStatusResponse, error) {
	var jsonData DmaOrderStatusResponse

	req := DmaOrderStatusRequest{
		UserId:       userId,
		VenueId:      venue,
		VenueOrderId: venueOrdId,
	}

	if bs, err := json.Marshal(&req); err != nil {
		return nil, err
	} else {
		resp, err := api.queryAdaptPrivateWithBytes("orders/status", bs)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resp, &jsonData)
		if err != nil {
			return nil, err
		}

		return &jsonData, nil
	}

}

// Function CancelOrderDMA cancels a DMA (direct market access) order.
func (api *Client) CancelOrderDMA(userId, venue, venueOrdId string) (*CancelDmaOrderResponse, error) {
	var jsonData CancelDmaOrderResponse

	req := CancelDmaOrderRequest{
		UserId:       userId,
		VenueId:      venue,
		VenueOrderId: venueOrdId,
	}

	if bs, err := json.Marshal(&req); err != nil {
		return nil, err
	} else {
		resp, err := api.queryAdaptPrivateWithBytes("orders/cancel", bs)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resp, &jsonData)
		if err != nil {
			return nil, err
		}

		return &jsonData, nil
	}

}

// Function GetConsolidatedOrderBookDMA gets current order book data across trading venues.
func (api *Client) GetConsolidatedOrderBookDMA(userId, asset, baseAsset string) (*DmaOrderBookResponse, error) {
	var jsonData DmaOrderBookResponse

	req := OrderBookRequest{
		UserId:    userId,
		Asset:     asset,
		BaseAsset: baseAsset,
	}

	if bs, err := json.Marshal(&req); err != nil {
		return nil, err
	} else {
		resp, err := api.queryAdaptPrivateWithBytes("data/real-time/order-book", bs)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resp, &jsonData)
		if err != nil {
			return nil, err
		}

		return &jsonData, nil
	}

}

// Function BalanceDMA provides the balance for a given asset at a given venue via the DMA API.
func (api *Client) BalanceDMA(userId, venue, assetId string) (*DmaBalanceResponse, error) {
	var jsonData DmaBalanceResponse

	req := DmaBalanceRequest{
		UserId:  userId,
		VenueId: venue,
		AssetId: assetId,
	}

	if bs, err := json.Marshal(&req); err != nil {
		return nil, err
	} else {
		resp, err := api.queryAdaptPrivateWithBytes("data/balance", bs)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resp, &jsonData)
		if err != nil {
			return nil, err
		}

		return &jsonData, nil
	}

}

// Function GetOrderStatus gets the current status of a Routefire (algorithm) order,
// to include amount filled and order open/closed flag.
func (api *Client) GetOrderStatus(userId string, orderId string) (*OrderStatusResponse, error) {
	var jsonData OrderStatusResponse
	params := map[string]interface{}{
		"user_id":  userId,
		"order_id": orderId,
	}

	resp, err := api.queryPrivate("orders/status", params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &jsonData)
	if err != nil {
		return nil, err
	}

	return &jsonData, nil
}

// Function CancelOrder cancels a Routefire (algorithm) order.
func (api *Client) CancelOrder(userId string, orderId string) (*OrderStatusResponse, error) {
	var jsonData OrderStatusResponse
	params := map[string]interface{}{
		"user_id":  userId,
		"order_id": orderId,
	}

	resp, err := api.queryPrivate("orders/cancel", params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &jsonData)
	if err != nil {
		return nil, err
	}

	return &jsonData, nil
}

// Function GetBalances gets the balances at each available trading venue for
// a given uid.
func (api *Client) GetBalances(uid, asset string) (map[string]string, error) {
	jsonData := map[string]string{}
	params := map[string]interface{}{
		"uid":   uid,
		"asset": asset,
	}

	resp, err := api.queryPrivate("data/balances", params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// Function GetOrderBookStats fetches key statistics from the order book for a
// prospective trade. The quantity provided is used to compute the "sweep cost,"
// or best theoretically available price given available liquidity.
func (api *Client) GetOrderBookStats(uid, buyAsset, sellAsset, quantity string) (*InquiryResponse, error) {
	var jsonData InquiryResponse
	params := map[string]interface{}{
		"uid":        uid,
		"buy_asset":  buyAsset,
		"sell_asset": sellAsset,
		"qty":        quantity,
	}

	resp, err := api.queryPrivate("data/inquire", params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &jsonData)
	if err != nil {
		return nil, err
	}

	return &jsonData, nil
}

// Function GetConsolidatedOrderBook fetches the order book for a given pair across exchanges.
func (api *Client) GetConsolidatedOrderBook(uid, buyAsset, sellAsset string) (*OrderBookResponse, error) {
	var jsonData OrderBookResponse
	params := map[string]interface{}{
		"uid":        uid,
		"buy_asset":  buyAsset,
		"sell_asset": sellAsset,
		"quantity":   "",
	}

	resp, err := api.queryPrivate("data/consolidated", params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &jsonData)
	if err != nil {
		return nil, err
	}

	return &jsonData, nil
}

func (api *Client) authenticate(uid, password string) (string, error) {
	var jsonData UserLoginResponse
	values := map[string]interface{}{
		"uid":      uid,
		"password": password,
	}

	resp, err := api.queryPublic("authenticate", values)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(resp, &jsonData)
	if err != nil {
		return "", err
	}

	if jsonData.Token != "" {
		api.accessToken = jsonData.Token
	}

	return jsonData.Token, nil
}

func (api *Client) queryPublic(command string, values map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", ApiUrl(), APIVersion, command)

	headers := map[string]string{"Content-Type": "application/json"}
	resp, err := api.doRequest(url, values, headers)

	return resp, err
}

func (api *Client) refreshToken() error {
	token, err := api.authenticate(api.username, api.password)
	if err != nil {
		return err
	}
	api.accessToken = token
	return nil
}

func (api *Client) refreshLoop(d time.Duration) {
	for {
		err := api.refreshToken()
		if err != nil {
			log.Printf("[RF] Failed to refresh auth token.\n")
		}
		time.Sleep(d)
	}
}

func (api *Client) queryAdaptPrivateWithBytes(command string, values []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", AdaptApiUrl(), APIVersion, command)
	log.Printf("Url: %s\n", url)

	// TODO: needs cleanup
	if len(api.accessToken) == 0 {
		if err0 := api.refreshToken(); err0 != nil {
			return nil, err0
		}
	}
	log.Printf("Auth token: %s\n", api.accessToken)
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", api.accessToken)}
	resp, err := api.doRequestBytes(url, values, headers)
	return resp, err
}

func (api *Client) queryPrivate(command string, values map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", ApiUrl(), APIVersion, command)

	// TODO: needs cleanup
	if len(api.accessToken) == 0 {
		if err0 := api.refreshToken(); err0 != nil {
			return nil, err0
		}
	}
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", api.accessToken)}
	resp, err := api.doRequest(url, values, headers)
	return resp, err
}

func (api *Client) QueryPrivate(command string, values map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", ApiUrl(), APIVersion, command)

	// TODO: needs cleanup
	if len(api.accessToken) == 0 {
		if err0 := api.refreshToken(); err0 != nil {
			return nil, err0
		}
	}
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", api.accessToken)}
	resp, err := api.doRequest(url, values, headers)
	return resp, err
}

func (api *Client) doRequest(reqURL string, params map[string]interface{}, headers map[string]string) ([]byte, error) {
	bytesParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(bytesParams))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", APIUserAgent)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := api.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (api *Client) doRequestBytes(reqURL string, bytesParams []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(bytesParams))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", APIUserAgent)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := api.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
