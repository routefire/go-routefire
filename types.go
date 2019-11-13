package routefire

type UserLoginResponse struct {
	Token string `json:"token"`
}

type MarketStatisticsResponse struct {
	UID       string `json:"uid"`
	BuyAsset  string `json:"buy_asset"`
	SellAsset string `json:"sell_asset"`
	Qty       string `json:"quantity"`
}

type SubmitOrderResponse struct {
	OrderId string `json:"order_id"`
}

type OrderStatusResponse struct {
	Status string `json:"status"`
	Filled string `json:"filled"`
}

type OrderBookEntry struct {
	Price    string
	Quantity string
	Venue    string
	Auction  bool
}

type OrderBook struct {
	Bids   []OrderBookEntry `json:"bids"`
	Offers []OrderBookEntry `json:"offers"`
}

type OrderBookResponse struct {
	Bids   []OrderBookEntry `json:"bids"`
	Offers []OrderBookEntry `json:"offers"`
}
