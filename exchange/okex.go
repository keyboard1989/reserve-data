package exchange

import (
	"fmt"
	"log"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Okex struct {
	interf       OkexInterface
	pairs        []common.TokenPair
	addresses    map[string]ethereum.Address
	exchangeInfo *common.ExchangeInfo
	fees         common.ExchangeFees
}

func (self *Okex) ID() common.ExchangeID {
	return common.ExchangeID("okex")
}

func (self *Okex) UpdateDepositAddress(token common.Token, address string) {
	self.addresses[token.ID] = ethereum.HexToAddress(address)
}

func (self *Okex) UpdatePrecisionLimit(pair common.TokenPair, exchangeInfo OkexInfo) {
	pairName := fmt.Sprintf("%s_%s", strings.ToLower(pair.Base.ID), strings.ToLower(pair.Quote.ID))
	for _, symbol := range exchangeInfo.Data {
		if symbol.Symbol == pairName {
			exchangePrecisionLimit := common.ExchangePrecisionLimit{}
			exchangePrecisionLimit.Precision.Amount = symbol.AmounPrecision
			exchangePrecisionLimit.Precision.Price = symbol.PricePrecision
			exchangePrecisionLimit.AmountLimit.Min = symbol.MinAmount
			self.exchangeInfo.Update(pair.PairID(), exchangePrecisionLimit)
			break
		}
	}
}

func (self *Okex) UpdatePairsPrecision() {
	exchangeInfo, err := self.interf.GetExchangeInfo()
	if err != nil {
		log.Printf("Get exchange info failed: %s\n", err)
	} else {
		for _, pair := range self.pairs {
			self.UpdatePrecisionLimit(pair, exchangeInfo)
		}
	}
}

func (self *Okex) GetInfo() (common.ExchangeInfo, error) {
	return *self.exchangeInfo, nil
}

func NewOkex(interf OkexInterface) *Okex {
	return &Okex{
		interf,
		[]common.TokenPair{
			common.MustCreateTokenPair("FUN", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
			common.MustCreateTokenPair("KNC", "ETH"),
			common.MustCreateTokenPair("LINK", "ETH"),
		},
		map[string]ethereum.Address{},
		common.NewExchangeInfo(),
		common.NewExchangeFee(
			common.TradingFee{
				"taker": 0.001,
				"maker": 0.001,
			},
			common.NewFundingFee(
				map[string]float32{},
				map[string]float32{},
			),
		),
	}
}
