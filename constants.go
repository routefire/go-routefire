package routefire

const (
	// Status codes
	StatusFilled          = "FILL"
	StatusError           = "ERROR"
	StatusPartiallyFilled = "PARTIAL_FILL"
	StatusOpen            = "OPEN"
	StatusCancelled       = "CANCEL"
	StatusComplete        = "COMPLETE"
	StatusExpired         = "EXPIRED"

	// Order side options
	SideBuy  = "BUY"
	SideSell = "SELL"
	SideCover = "COVER"
	SideShort = "SHORT"

	// Venue IDs
	CoinbasePro = "GDAX"
	Gemini      = "GEMINI"
	Binance     = "BINANCE"
	Bittrex     = "BITTREX"
	Kraken      = "KRAKEN"
	Bitfinex    = "BITFINEX"
	Poloniex    = "POLONIEX"

	// Fiat currency codes
	Usd = "usd"
	Eur = "eur"
	Gbp = "gbp"

	// Cryptocurrency codes (stablecoins)
	Usdt = "usdt"
	Usdc = "usdc"
	Tusd = "tusd"
	Gusd = "gusd"
	Dai  = "dai"
	Pax  = "pax"

	// Cryptocurrency codes (other)
	Btc = "btc"
	Bch = "bch"
	Eth = "eth"
	Ltc = "ltc"
	Xrp = "xrp"
	Xlm = "xlm"
	Zrx = "zrx"
)
