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

type TokenPairID string

func NewTokenPairID(base, quote string) TokenPairID {
	return TokenPairID(fmt.Sprintf("%s-%s", base, quote))
}

type ExchangeID string

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

type RawBalance big.Int

func (self *RawBalance) ToFloat(decimal int64) float64 {
	f := new(big.Float).SetInt((*big.Int)(self))
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(decimal), nil,
	))
	res := new(big.Float).Quo(f, power)
	result, _ := res.Float64()
	return result
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
