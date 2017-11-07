package exchange

// map of token pair to map of asks/bids to array of [rate, amount]
type Liqresp map[string]map[string][][]float64

type Liqinfo struct {
	Success int `json:"success"`
	Return  map[string]map[string]float64
	Error   string `json:"error"`
}

type Liqwithdraw struct {
	Success int `json:"success"`
	Return  map[string]interface{}
	Error   string `json:"error"`
}

type Liqtrade struct {
	Success int `json:"success"`
	Return  struct {
		Done      float64 `json:"received"`
		Remaining float64 `json:"remains"`
		OrderID   int64   `json:"order_id"`
	}
	Error string `json:"error"`
}
