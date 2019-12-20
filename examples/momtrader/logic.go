package main

import (
	"errors"
	"fmt"
	"github.com/routefire/go-routefire"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

type position struct {
	asset string
	size  string
	price string
}

type momentumParams struct {
	gainerPeriods int
	stdDevPeriods int
	minStdDevBuy  float64
}

type order struct {
	order      *routefire.PlaceDmaOrderResponse
	isFilled   bool
	isComplete bool
}

type MomentumTrader struct {
	UserId        string
	Assets        []string
	BaseAsset     string
	Alpha         float64
	Positions     []position
	Orders        []*order
	BidHistory    map[string][]float64
	AskHistory    map[string][]float64
	LastOrderBook map[string]*routefire.DmaOrderBookResponse
	RfClient      *routefire.Client
	Capital       float64
	params        *momentumParams
	lock          *sync.Mutex
}

func NewMomentumTrader(uid string, rfClient *routefire.Client, assets []string, baseAsset string, capital, alpha float64) *MomentumTrader {
	m1 := map[string][]float64{}
	m2 := map[string][]float64{}
	obm := map[string]*routefire.DmaOrderBookResponse{}
	return &MomentumTrader{
		UserId:        uid,
		Assets:        assets,
		BaseAsset:     baseAsset,
		Alpha:         alpha,
		Positions:     nil,
		Orders:        nil,
		LastOrderBook: obm,
		BidHistory:    m1,
		AskHistory:    m2,
		Capital:       capital,
		RfClient:      rfClient,
		params: &momentumParams{
			gainerPeriods: 5,
			stdDevPeriods: 15,
			minStdDevBuy:  alpha,
		},
		lock: &sync.Mutex{},
	}
}

func (m *MomentumTrader) RunLoop(dur time.Duration) {
	for {
		err := m.Trade()
		if err != nil {
			m.log("Trade() error: %s", err.Error())
		}
		time.Sleep(dur)
	}
}
func (m *MomentumTrader) collectData() error {
	for _, asset := range m.Assets {
		ob, err := m.RfClient.GetConsolidatedOrderBookDMA(m.UserId, asset, m.BaseAsset)
		if err != nil {
			return err
		} else {
			m.LastOrderBook[asset] = ob
		}
		bestOff := ob.Data.Offers[0]
		bestBid := ob.Data.Bids[len(ob.Data.Bids)-1]

		bestOffF, err := strconv.ParseFloat(bestOff.Price, 64)
		bestBidF, err := strconv.ParseFloat(bestBid.Price, 64)

		m.BidHistory[asset] = append(m.BidHistory[asset], bestBidF)
		m.AskHistory[asset] = append(m.AskHistory[asset], bestOffF)
	}
	return nil
}

func (m *MomentumTrader) gainOver(asset string, nPeriods int) (float64, error) {
	if arr, ok := m.AskHistory[asset]; ok {
		if len(arr) < nPeriods {
			return 0.0, errors.New("InsufficientData")
		} else {
			p0 := arr[len(arr)-nPeriods]
			p1 := arr[len(arr)-1]
			gain := (p1 / p0) - 1.0
			return gain, nil
		}
	}
	return 0.0, errors.New("InsufficientData")
}

func (m *MomentumTrader) biggestGainerOver(nPeriods int) (string, float64, error) {
	bestYet := -999999.99
	bestYetAsst := ""
	for _, asset := range m.Assets {
		g, err := m.gainOver(asset, nPeriods)
		if err != nil {
			continue
		}
		if g > bestYet {
			bestYet = g
			bestYetAsst = asset
		}
		m.log("\tGains: %s gained %f...", asset, g)
	}
	if len(bestYetAsst) == 0 {
		return "", 0.0, errors.New("InsufficientData")
	}
	return bestYetAsst, bestYet, nil
}

func (m *MomentumTrader) stdDevs(asset string, nPeriods int) (float64, error) {
	if arr0, ok := m.AskHistory[asset]; ok && len(arr0) >= nPeriods {
		arr := arr0[len(arr0)-nPeriods:]
		sigma := stdDev(arr)
		if math.Abs(sigma) < 0.00000001 {
			return 0.0, nil
		}
		mu := mean(arr)
		last := arr[len(arr)-1]
		sigmas := (last - mu) / sigma
		m.log("Std devs for %s: last %f - average %f / sigma %f = %f", asset, last, mu, sigma, sigmas)
		return sigmas, nil
	}
	return 0.0, errors.New("InsufficientData")
}

func (m *MomentumTrader) Trade() error {
	m.lock.Lock()

	err := m.collectData()
	if err != nil {
		m.log("Data collection error: %s", err.Error())
		m.lock.Unlock()
		return err
	}

	winner, _, err := m.biggestGainerOver(5)
	if err != nil {
		m.log("Math/data error: %s", err.Error())
		m.lock.Unlock()
		return err
	}

	wSds, err := m.stdDevs(winner, 10)
	if err != nil {
		m.log("Math/data error [2]: %s", err.Error())
		m.lock.Unlock()
		return err
	}
	m.log("Winner is %s - at %f SDs", winner, wSds)

	if wSds >= m.params.minStdDevBuy {
		m.log("Buy %s!", winner)
	} else {
		m.log("Std dev not met for: %s", winner)
		m.lock.Unlock()
		return err
	}
	m.lock.Unlock()

	if m.waitingOnExecution() {
		m.log("Skipping iteration, waiting on execution...")
		return nil
	}

	bb := m.LastOrderBook[winner].Data.Bids[len(m.LastOrderBook[winner].Data.Bids)-1]
	bo := m.LastOrderBook[winner].Data.Offers[0]

	px := bo.Price
	size := m.amountToTradeAt(px)
	m.log("Intended positions: %s %s @ %s (%s)", size, winner, px, bb.Venue)

	if len(m.Positions) > 0 && m.Positions[0].asset != winner {
		// Sell the old position
		curPos := m.Positions[0]
		//m.log("Need to sell %s %s @ %s (%s)", curPos.asset, curPos.size, bb.Price, bo.Venue)
		m.log("EXIT %s %s @ %s (%s)", curPos.size, curPos.asset, bb.Price, curPos)
		err := m.doTrade(curPos.asset, bo.Venue, curPos.size, bb.Price, false)
		if err != nil || DevelopmentExecutionSafety {
			return err
		}

		// Do the buy
		m.log("ENTER %s %s @ %s (%s)", size, winner, px, bb.Venue)
		err = m.doTrade(winner, bb.Venue, size, px, true)
		if err != nil {
			return err
		}
	} else if len(m.Positions) == 0 {
		// Do the buy
		m.log("ENTER[0] %s %s @ %s (%s)", size, winner, px, bb.Venue)
		err = m.doTrade(winner, bb.Venue, size, px, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MomentumTrader) priceAtVenue(asst, venu string, isBid bool) string {
	//
}

func (m *MomentumTrader) amountToTradeAt(px string) string {
	pxF, err := strconv.ParseFloat(px, 64)
	if err != nil {
		return "0"
	}
	amt := m.Capital / pxF
	amtS := fmt.Sprintf("%.4f", amt)
	return amtS
}

func (m *MomentumTrader) doTrade(asset, venu, qty, px string, isBuy bool) error {
	ch := make(chan error, 1)

	res, err := SubmitAndWait(m.UserId, m.RfClient, isBuy, asset, m.BaseAsset, venu, qty, px, ch, 3*time.Second)
	if err != nil {
		return err
	}

	ord := &order{
		order:      res,
		isFilled:   false,
		isComplete: false,
	}
	m.Orders = append(m.Orders, ord)

	go func(errCh chan error, asst, size, price string, m0 *MomentumTrader, ordr *order) {
		mustCancel := false
		select {
		case err := <-errCh:
			if err == nil {
				if isBuy {
					// buy...
					m.Positions = append(m.Positions, position{
						asset: asset,
						size:  qty,
						price: px,
					})
				} else {
					// sell...
					m0.removePosition(asst, size, price)
				}
				m0.setOrderComplete(ordr.order.VenueId, ordr.order.VenueOrderId, true)
			} else {
				m0.log("CRITICAL - Order failed: %s", err.Error())
				m0.setOrderComplete(ordr.order.VenueId, ordr.order.VenueOrderId, false)
				mustCancel = true
			}
		case <-time.NewTicker(60 * time.Second).C:
			m0.log("CRITICAL - Timed out waiting for trade to finish...")
			m0.setOrderComplete(ordr.order.VenueId, ordr.order.VenueOrderId, false)
			mustCancel = true
		}
		if mustCancel {
			m.log("Canceling unsuccessful order %s...")
			_, err := m0.RfClient.CancelOrderDMA(m0.UserId, ord.order.VenueId, ord.order.VenueOrderId)
			if err != nil {
				panic(err)
			}
		}
	}(ch, asset, qty, px, m, ord)

	return nil
}

func (m *MomentumTrader) waitingOnExecution() bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, x := range m.Orders {
		if !x.isComplete {
			return true
		}
	}

	return false
}

func (m *MomentumTrader) removePosition(asset, qty, px string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var newPos []position
	for _, x := range m.Positions {
		if !(x.asset == asset && x.price == px && x.size == qty) {
			newPos = append(newPos, x)
		}
	}
	m.Positions = newPos
}

func (m *MomentumTrader) setOrderComplete(venu, venuOrdId string, filled bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for i, x := range m.Orders {
		if x.order.VenueOrderId == venuOrdId && x.order.VenueId == venu {
			m.Orders[i].isComplete = true
			if filled {
				m.Orders[i].isFilled = true
			}
		}
	}
}

func (m *MomentumTrader) log(fmtStr string, argArr ...interface{}) {
	fs := "MT> " + fmtStr + "\n"
	log.Printf(fs, argArr...)
}
