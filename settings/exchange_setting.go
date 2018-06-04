package settings

// ExchangeName is the name of exchanges of which core will use to rebalance.
//go:generate stringer -type=ExchangeName
type ExchangeName int

const (
	Binance ExchangeName = iota
	Bittrex
	Huobi
	StableExchange
)

var exchangeNameValue = map[string]ExchangeName{
	"binance":         Binance,
	"bittrex":         Bittrex,
	"huobi":           Huobi,
	"stable_exchange": StableExchange,
}

func ExchangTypeValues() map[string]ExchangeName {
	return exchangeNameValue
}

type ExchangeSetting struct {
	Storage ExchangeStorage
}

func NewExchangeSetting(exchangeStorage ExchangeStorage) (*ExchangeSetting, error) {
	return &ExchangeSetting{exchangeStorage}, nil
}

func (setting *Settings) LoadFeeFromFile(path string) error {
	return nil
}

func (setting *Settings) LoadDepositAddressFromFile(path string) error {
	return nil
}

func (setting *Settings) LoadMinDepositAddressFromFile(path string) error {
	return nil
}
