package exchange

type Bitfresp struct {
	Asks []map[string]string `json:"asks"`
	Bids []map[string]string `json:"bids"`
}

type Bitfinfo struct {
	Success int `json:"success"`
	Return  map[string]map[string]float64
	Error   string `json:"error"`
}

type BitExchangeInfo struct {
	Pairs []struct {
		Pair           string `json:"pair"`
		PricePrecision int    `json:"price_precision"`
	}
}

type Bitfwithdraw struct {
	Success int `json:"success"`
	Return  map[string]interface{}
	Error   string `json:"error"`
}

type Bitftrade struct {
	Success int `json:"success"`
	Return  struct {
		Done      float64 `json:"received"`
		Remaining float64 `json:"remains"`
		OrderID   int64   `json:"order_id"`
	}
	Error string `json:"error"`
}
