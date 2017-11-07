package exchange

import (
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Binance struct {
	signer    Signer
	endpoint  BinanceEndpoint
	pairs     []common.TokenPair
	addresses map[string]ethereum.Address
}

func (self Binance) ID() common.ExchangeID {
	return common.ExchangeID("binance")
}

func (self Binance) Name() string {
	return "binance"
}

func (self *Binance) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Binance) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses[token.ID]
	return addr, supported
}

func (self *Binance) TokenPairs() []common.TokenPair {
	return self.pairs
}


func (self *Binance) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64) (done float64, remaining float64, finished bool, err error) {
	return self.endpoint.Trade(
		self.signer.GetBinanceKey(),
		tradeType, base, quote, rate, amount, self.signer)
}

func (self *Binance) Withdraw(token common.Token, amount *big.Int, address ethereum.Address) error {
	return self.endpoint.Withdraw(
		self.signer.GetBinanceKey(),
		token, amount, address, self.signer)
}

func (self Binance) FetchPriceData() (map[common.TokenPairID]common.ExchangePrice, error) {
	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs := self.pairs
	for _, pair := range pairs {
		wait.Add(1)
		go self.endpoint.FetchOnePairData(&wait, pair, &data)
	}
	wait.Wait()
	result := map[common.TokenPairID]common.ExchangePrice{}
	data.Range(func(key, value interface{}) bool {
		result[key.(common.TokenPairID)] = value.(common.ExchangePrice)
		return true
	})
	// fmt.Printf("result: %v\n", result)
	return result, nil
}

func (self Binance) FetchEBalanceData() (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.GetTimestamp()
	result.Valid = true
	//response, err := self.endpoint.GetInfo(
	//	self.signer.getBinanceKey(),
	//	self.signer,
	//)
	result.ReturnTime = common.GetTimestamp()
	//if err != nil {
	//	result.Valid = false
	//	result.Error = error.Error()
	//}
	return result, nil
}

func NewBinance(signer Signer, endpoint BinanceEndpoint) *Binance {
	return &Binance{
		signer,
		endpoint,
		[]common.TokenPair{
			common.MustCreateTokenPair("FUN", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
			common.MustCreateTokenPair("KNC", "ETH"),
			common.MustCreateTokenPair("LINK", "ETH"),
		},
		map[string]ethereum.Address{
			"ETH": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"OMG": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"DGD": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"CVC": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"MCO": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"GNT": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"ADX": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"EOS": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"PAY": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"BAT": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"KNC": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
		},
	}
}
