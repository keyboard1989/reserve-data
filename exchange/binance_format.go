package exchange

type Binaresp struct {
	LastUpdatedId int64      `json:"lastUpdateId"`
	Code          int        `json:"code"`
	Msg           string     `json:"msg"`
	Bids          [][]string `json:"bids"`
	Asks          [][]string `json:"asks"`
}

type Binainfo struct {
	MakerCommission  int64 `json:"makerCommission"`
	TakerCommission  int64 `json:"takerCommission"`
	BuyerCommission  int64 `json:"buyerCommission"`
	SellerCommission int64 `json:"sellerCommission"`
	CanTrade         bool  `json:"canTrade"`
	CanWithdraw      bool  `json:"canWithdraw"`
	CanDeposit       bool  `json:"canDeposit"`
	Balances []struct {
		Asset  string  `json:"asset"`
		Free   float64 `json:"free"`
		Locked float64 `json:"locked`
	} `json:"balances`
}

type Binatrade struct{}

type Binawithdraw struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
}
