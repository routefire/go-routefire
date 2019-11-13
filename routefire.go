package routefireApi

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	APIURL       = "https://routefire.io/api"
	APIVersion   = "v1"
	APIUserAgent = "RouteFire API client agent"
)

type RouteFireAPI struct {
	username    string
	password    string
	accessToken string
	client      *http.Client
}

func New(uid, password string) *RouteFireAPI {
	return &RouteFireAPI{uid, password, "", http.DefaultClient}
}

func (api *RouteFireAPI) SubmitOrder(userId string, buyAsset string, sellAsset string, quantity string, price string, algo string, algoParams map[string]string) (*SubmitOrderResponse, error) {
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
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	return &jsonData, nil
}

func (api *RouteFireAPI) GetOrderStatus(userId string, orderId string) (*OrderStatusResponse, error) {
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
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	return &jsonData, nil
}

func (api *RouteFireAPI) CancelOrder(userId string, orderId string) (*OrderStatusResponse, error) {
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
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	return &jsonData, nil
}

func (api *RouteFireAPI) GetConsolidatedOrderBook(uid, buyAsset, sellAsset, quantity string) (*OrderBookResponse, error) {
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
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	return &jsonData, nil
}

func (api *RouteFireAPI) authenticate(uid, password string) (string, error) {
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

func (api *RouteFireAPI) queryPublic(command string, values map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", APIURL, APIVersion, command)

	headers := map[string]string{"Content-Type": "application/json"}
	resp, err := api.doRequest(url, values, headers)

	return resp, err
}

func (api *RouteFireAPI) queryPrivate(command string, values map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", APIURL, APIVersion, command)
	if api.accessToken != "" {
		headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", api.accessToken)}
		resp, err := api.doRequest(url, values, headers)
		return resp, err
	}

	token, err := api.authenticate(api.username, api.password)
	if err != nil {
		return nil, fmt.Errorf("Could not get JWT token (%s)", err.Error())
	}

	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	resp, err := api.doRequest(url, values, headers)
	return resp, err
}

func (api *RouteFireAPI) doRequest(reqURL string, params map[string]interface{}, headers map[string]string) ([]byte, error) {
	bytesParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! (%s)", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(bytesParams))
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #1 (%s)", err.Error())
	}

	req.Header.Add("User-Agent", APIUserAgent)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := api.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #2 (%s)", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #3 (%s)", err.Error())
	}

	return body, nil
}
