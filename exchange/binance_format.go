package exchange

type Binaresp struct {
	LastUpdatedId int64      `json:"lastUpdateId"`
	Code          int        `json:"code"`
	Msg           string     `json:"msg"`
	Bids          [][]string `json:"bids"`
	Asks          [][]string `json:"asks"`
}

type Binainfo struct {
	Success int `json:"success"`
	Return  map[string]map[string]float64
	Error   string `json:"error"`
}