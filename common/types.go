package common

import (
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type Version int64
type Timestamp string

func GetTimestamp() Timestamp {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return Timestamp(strconv.Itoa(int(timestamp)))
}

func TimeToTimepoint(t time.Time) uint64 {
	timestamp := t.UnixNano() / int64(time.Millisecond)
	return uint64(timestamp)
}

type TokenPairID string

func NewTokenPairID(base, quote string) TokenPairID {
	return TokenPairID(fmt.Sprintf("%s-%s", base, quote))
}

type ExchangeID string

type ActivityRecord struct {
	ID        int64
	Timestamp Timestamp
	Action    string
	Params    map[string]interface{}
	Result    interface{}
}

type PriceEntry struct {
	Quantity float64
	Rate     float64
}

type AllPriceResponse struct {
	Version    Version
	Timestamp  Timestamp
	ReturnTime Timestamp
	Data       map[TokenPairID]OnePrice
}

type OnePriceResponse struct {
	Version    Version
	Timestamp  Timestamp
	ReturnTime Timestamp
	Data       OnePrice
}

type OnePrice map[ExchangeID]ExchangePrice

type ExchangePrice struct {
	Valid      bool
	Error      string
	Timestamp  Timestamp
	BuyPrices  []PriceEntry
	SellPrices []PriceEntry
	ReturnTime Timestamp
}

func BigToFloat(b *big.Int, decimal int64) float64 {
	f := new(big.Float).SetInt(b)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(decimal), nil,
	))
	res := new(big.Float).Quo(f, power)
	result, _ := res.Float64()
	return result
}

type RawBalance big.Int

func (self *RawBalance) ToFloat(decimal int64) float64 {
	return BigToFloat((*big.Int)(self), decimal)
}

type BalanceEntry struct {
	Valid      bool
	Error      string
	Timestamp  Timestamp
	ReturnTime Timestamp
	Balance    RawBalance
}

func (self BalanceEntry) ToBalanceResponse(decimal int64) BalanceResponse {
	return BalanceResponse{
		Valid:      self.Valid,
		Error:      self.Error,
		Timestamp:  self.Timestamp,
		ReturnTime: self.ReturnTime,
		Balance:    self.Balance.ToFloat(decimal),
	}
}

type BalanceResponse struct {
	Valid      bool
	Error      string
	Timestamp  Timestamp
	ReturnTime Timestamp
	Balance    float64
}

type AllBalanceResponse struct {
	Version    Version
	Timestamp  Timestamp
	ReturnTime Timestamp
	Data       map[string]BalanceResponse
}

type EBalanceEntry struct {
	Valid      bool
	Error      string
	Timestamp  Timestamp
	ReturnTime Timestamp
	Balance    map[string]float64
}

type AllEBalanceResponse struct {
	Version    Version
	Timestamp  Timestamp
	ReturnTime Timestamp
	Data       map[ExchangeID]EBalanceEntry
}

type RateEntry struct {
	Rate        *big.Int
	ExpiryBlock *big.Int
	Balance     *big.Int
}

type RateResponse struct {
	Valid       bool
	Error       string
	Timestamp   Timestamp
	ReturnTime  Timestamp
	Rate        float64
	ExpiryBlock int64
	Balance     float64
}

type AllRateEntry struct {
	Valid      bool
	Error      string
	Timestamp  Timestamp
	ReturnTime Timestamp
	Data       map[TokenPairID]RateEntry
}

type AllRateResponse struct {
	Version    Version
	Valid      bool
	Error      string
	Timestamp  Timestamp
	ReturnTime Timestamp
	Data       map[TokenPairID]RateResponse
}
