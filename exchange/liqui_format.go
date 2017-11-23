package exchange

// map of token pair to map of asks/bids to array of [rate, amount]
type Liqresp map[string]map[string][][]float64

type Liqinfo struct {
	Success int                           `json:"success"`
	Return  map[string]map[string]float64 `json:"return"`
	Error   string                        `json:"error"`
}

type Liqwithdraw struct {
	Success int                    `json:"success"`
	Return  map[string]interface{} `json:"return"`
	Error   string                 `json:"error"`
}

type Liqtrade struct {
	Success int `json:"success"`
	Return  struct {
		Done      float64 `json:"received"`
		Remaining float64 `json:"remains"`
		OrderID   uint64  `json:"order_id"`
	} `json:"return"`
	Error string `json:"error"`
}

type Liqorders struct {
	Success int `json:"success"`
	Return  map[string]struct {
		Pair      string  `json:"pair"`
		Type      string  `json:"type"`
		Amount    float64 `json:"amount"`
		Rate      float64 `json:"rate"`
		Timestamp uint64  `json:"timestamp_created"`
	} `json:"return"`
	Error string `json:"error"`
}

type Liqcancel struct {
	Success int `json:"success"`
	Return  struct {
		OrderID uint64             `json:"order_id"`
		Funds   map[string]float64 `json:"funds"`
	} `json:"return"`
	Error string `json:"error"`
}
