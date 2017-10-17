package market

type Price struct {
	Quantity float64
	Rate     float64
}

type ExchangeData struct {
	Valid      bool
	Timestamp  string
	BuyPrices  []Price
	SellPrices []Price
}

type OnePairData map[string]ExchangeData

type PriceData struct {
	Version      int64
	Timestamp    string
	ExchangeData OnePairData
}

type AllPriceData struct {
	Version     int64
	Timestamp   string
	AllPairData map[string]OnePairData
}

func NewAllPriceData() *AllPriceData {
	return &AllPriceData{
		0, "", map[string]OnePairData{},
	}
}

var GetPrice func(string, string) (PriceData, error)
var GetAllPrice func() (AllPriceData, error)
