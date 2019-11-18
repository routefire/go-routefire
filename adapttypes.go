package routefire

type DmaError struct {
	Message string `json:"error"`
}

type PlaceDmaOrderRequest struct {
	UserId      string            `json:"user_id"`
	VenueId     string            `json:"venue"`
	Side        string            `json:"side"`
	TradedAsset string            `json:"traded_asset"`
	BaseAsset   string            `json:"base_asset"`
	Quantity    string            `json:"quantity"`
	Price       string            `json:"price"`
	OrderParams map[string]string `json:"order_params"`
}

type PlaceDmaOrderResponse struct {
	VenueId      string     `json:"venue"`
	VenueOrderId string     `json:"venue_order_id"`
	Errors       []DmaError `json:"errors"`
}

type CancelDmaOrderRequest struct {
	UserId       string `json:"user_id"`
	VenueId      string `json:"venue"`
	VenueOrderId string `json:"venue_order_id"`
}

type CancelDmaOrderResponse struct {
	VenueId      string     `json:"venue"`
	VenueOrderId string     `json:"venue_order_id"`
	Errors       []DmaError `json:"errors"`
}

type DmaOrderStatusRequest struct {
	UserId       string `json:"user_id"`
	VenueId      string `json:"venue"`
	VenueOrderId string `json:"venue_order_id"`
}

type DmaOrderStatusResponse struct {
	VenueId      string     `json:"venue"`
	VenueOrderId string     `json:"venue_order_id"`
	Status       string     `json:"status"`
	FilledAmount string     `json:"filled"`
	Errors       []DmaError `json:"errors"`
}

type OrderBookRequest struct {
	UserId    string `json:"user_id"`
	Asset     string `json:"asset"`
	BaseAsset string `json:"base_asset"`
}

type DmaOrderBookResponse struct {
	Data   DmaOrderBook `json:"data"`
	Errors []DmaError   `json:"errors"`
}

type DmaOrderBook struct {
	Bids   []DmaOrderBookEntry `json:"bids"`
	Offers []DmaOrderBookEntry `json:"offers"`
}

type DmaOrderBookEntry struct {
	Amount    string `json:"quantity"`
	Price     string `json:"price"`
	Venue     string `json:"venue"`
	BuyAsset  string `json:"buy_asset"`
	SellAsset string `json:"sell_asset"`
}

type DmaBalanceRequest struct {
	UserId  string `json:"user_id"`
	VenueId string `json:"venue"`
	AssetId string `json:"asset"`
}

type DmaBalanceResponse struct {
	VenueId string     `json:"venue"`
	Asset   string     `json:"asset"`
	Amount  string     `json:"amount"`
	Errors  []DmaError `json:"errors"`
}
