# go-routefire: a native Go SDK for the Routefire API

Trade cryptocurrencies and digital assets across all major exchanges with a single, 
stable, unified API.

## Setup

### Importing 

Simply import as normal: 

```go
import (
    "github.com/routefire/go-routefire"
)
``` 

`vgo` modules are enabled, but legacy-style compilation is supported as well.

### Getting an account

You will at minimum need a free Routefire account, which can be obtained for free at
 the [Routefire web site](https://routefire.io/signup?q=oss).
 
## Usage

The simplest way to use the API is using username/password authentication. To do this,
simply call the `New` function:

```go
client := routefire.New(username, password)
```

### Direct market access (DMA) orders

The DMA API provides low-level access to the connectivity layer in Routefire Core. 
Therefore, DMA orders specify precisely the venue and price level at which to place 
a trade, instead of using an algorithm to decide the optimal way to enter the order.
(A standard free Routefire account is a DMA account.)

The DMA API is available via the methods ending in `*DMA`: 
- `SubmitOrderDMA`: submit an order to a trading venue
- `OrderStatusDMA`: get order status from a trading venue
- `CancelOrderDMA`: cancel a given order at a trading venue
- `GetConsolidatedOrderBookDMA`: get consolidated order book across trading venues 
- `BalanceDMA`: get balance for a given asset at a given venue 

### Routefire (algorithmic) orders

To submit orders that are worked by Routefire algorithms, a different set of methods
is used from DMA (direct market access) modules. The unit tests in `routefire_test.go`
demonstrate how these orders are parameterized: each algorithm has a unique set of
parameters that it accepts; the parameters used in the unit test are the most 
commonly used.

To submit an order, call `SubmitOrder`:

```go
params := map[string]string{
	"target_seconds": "100",
	"backfill":       "1.0",
	"aggression":     "0.0",
}

resp, err := client.SubmitOrder(uid, "btc", "usd", "0.003", "", "rfxw", params)
```

This submits an algorithmic order to buy 0.003 BTC using USD via the RFXW trading
algorithm. RFXW is instructed to target 100 seconds to fill the order, and the `1.0`
value given to `backfill` indicates that liquidity can be taken to avoid falling
behind schedule. Note that no price is provided or needed for algorithmic orders.

The order ID for the new order (assuming submission was successful) will be contained in
the `OrderId` field of `resp`. This ID can be used in subsequent calls to either check
the status of or cancel the order. For example:

```go
status, err := client.GetOrderStatus(uid, resp.OrderId)
```

Or:

```go
status, err := client.CancelOrder(uid, resp.OrderId)
```

## Schema

### Handling numbers

All functions accept string values for prices and quantities (to preserve numerical
precision). 

### Important constants

String identifiers are used to specify assets, trading venues, and side.
These constants are provided in `constants.go`. Most importantly, there are:
 
- *Assets*: e.g. `Usd`, `Btc` 
- *Trading venues*: e.g. `CoinbasePro`, `Binance`
- *Side*: `SideBuy`, `SideSell`, `SideShort`, `SideCover`

## Examples

Examples are provided in the `examples` directory of this repository:

- `data`: demonstrates how to fetch order book data and account balances using 
  both the Core and DMA APIs.
- `dma`: demonstrates how to use the DMA API to submit an order, watch it fill,
  and expire it if it doesn't fill within five minutes.
