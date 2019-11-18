# go-routefire: a native Go SDK for the Routefire API

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
 the [Routefire web site](https://routefire.io).
 
## Usage

The simplest way to use the API is using username/password authentication. To do this,
simply call the `New` function:

```go
client := routefire.New(username, password)
```

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

resp, err := client.SubmitOrder(uid, "btc", "usd", "0.003", "10000.0", "rfxw", params)
```

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

### Direct market access (DMA) orders

The DMA API provides low-level access to the connectivity layer in Routefire Core. 
Therefore, DMA orders specify precisely the venue and price level at which to place 
a trade, instead of using an algorithm to decide the optimal way to enter the order.

The DMA API is available via the methods ending in `*DMA`: `SubmitOrderDMA`, 
`OrderStatusDMA`, `CancelOrderDMA`, `GetConsolidatedOrderBookDMA`, and
`BalanceDMA`. 

In general, the usage of these functions is straightforward and uses 
venue-generated order IDs exclusively. Relevant asset IDs, venue IDs, and
other useful constants are provided in `costants.go` (e.g. `Usd`, `Btc`, 
`CoinbasePro`, etc.).

