package exchange

// map of token pair to map of asks/bids to array of [rate, amount]
type Bittresp struct {
	Success bool                            `json:"success"`
	Msg     string                          `json:"message"`
	Result  map[string][]map[string]float64 `json:"result"`
}

type Bittinfo struct {
	Success int    `json:"success"`
	Return  map[string]map[string]float64
	Error   string `json:"error"`
}

type Bittwithdraw struct {
	Success int    `json:"success"`
	Return  map[string]interface{}
	Error   string `json:"message"`
}

type Bitttraderesult struct {
	Done      float64 `json:"received"`
	Remaining float64 `json:"remains"`
	OrderID   int64   `json:"uuid"`
}

type Bitttrade struct {
	Success int             `json:"success"`
	Return  Bitttraderesult `json:"result"`
	Error   string          `json:"message"`
}
