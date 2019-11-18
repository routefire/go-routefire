package routefire

// Type UserLoginResponse holds a JWT token returned by an authentication endpoint.
type UserLoginResponse struct {
	Token string `json:"token"`
}

// Type MarketStatisticsResponse holds order book statistics requested from
// the `/inquire` endpoint.
type MarketStatisticsResponse struct {
	UID       string `json:"uid"`
	BuyAsset  string `json:"buy_asset"`
	SellAsset string `json:"sell_asset"`
	Qty       string `json:"quantity"`
}

// Type SubmitOrderResponse object provides an order ID for a new order, or an
// empty string if the order could not be submitted.
type SubmitOrderResponse struct {
	OrderId string `json:"order_id"`
}

// Type OrderStatusResponse object provides the current status and filled amount
// for the requested order.
type OrderStatusResponse struct {
	Status string `json:"status"`
	Filled string `json:"filled"`
}

// Type OrderBookEntry represents a single order book line item (a row in the L2 depth data).
type OrderBookEntry struct {
	Price    string
	Quantity string
	Venue    string
	Auction  bool
}

// Type OrderBook provides order book data -- a bid and ask side.
type OrderBook struct {
	Bids   []OrderBookEntry `json:"bids"`
	Offers []OrderBookEntry `json:"offers"`
}

// Type OrderBookResponse provides order book data -- a bid and ask side.
type OrderBookResponse struct {
	Bids   []OrderBookEntry `json:"bids"`
	Offers []OrderBookEntry `json:"offers"`
}

// Type InquiryResponse holds order book statistics requested from the `/inquire` endpoint.
type InquiryResponse struct {
	IsoCost            float64            `json:"iso_cost"`
	TopPrices          map[string]float64 `json:"top_of_book"`
	TopPriceChanges    map[string]float64 `json:"top_of_book_changes"`
	TopPriceChangesPct map[string]float64 `json:"top_of_book_changes_pct"`
}


