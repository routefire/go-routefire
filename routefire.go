package routefire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	APIURL       = "https://routefire.io/api"
	APIVersion   = "v1"
	APIUserAgent = "RouteFire API client agent"
	AuthInterval = 100
)

var webHttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 10,
	},
	Timeout: time.Second * 20,
}

type Client struct {
	username    string
	password    string
	accessToken string
	client      *http.Client
}

func New(uid, password string) *Client {
	z := &Client{uid, password, "", webHttpClient}
	go z.refreshLoop(AuthInterval*time.Second)
	return z
}

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

func (api *Client) GetConsolidatedOrderBook(uid, buyAsset, sellAsset, quantity string) (*OrderBookResponse, error) {
	var jsonData OrderBookResponse
	params := map[string]interface{}{
		"uid":        uid,
		"buy_asset":  buyAsset,
		"sell_asset": sellAsset,
		"quantity":   quantity,
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
	url := fmt.Sprintf("%s/%s/%s", APIURL, APIVersion, command)

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

func (api *Client) queryPrivate(command string, values map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", APIURL, APIVersion, command)

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
