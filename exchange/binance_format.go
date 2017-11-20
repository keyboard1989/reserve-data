package exchange

type Binaresp struct {
	LastUpdatedId int64      `json:"lastUpdateId"`
	Code          int        `json:"code"`
	Msg           string     `json:"msg"`
	Bids          [][]string `json:"bids"`
	Asks          [][]string `json:"asks"`
}

type Binainfo struct {
	Code             int    `json:"code"`
	Msg              string `json:"msg"`
	MakerCommission  int64  `json:"makerCommission"`
	TakerCommission  int64  `json:"takerCommission"`
	BuyerCommission  int64  `json:"buyerCommission"`
	SellerCommission int64  `json:"sellerCommission"`
	CanTrade         bool   `json:"canTrade"`
	CanWithdraw      bool   `json:"canWithdraw"`
	CanDeposit       bool   `json:"canDeposit"`
	Balances         []struct {
		Asset  string  `json:"asset"`
		Free   float64 `json:"free"`
		Locked float64 `json:"locked`
	} `json:"balances`
}

type Binatrade struct {
	Symbol        string `json:"symbol"`
	OrderID       uint64 `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  uint64 `json:"transactTime"`
}

type Binawithdraw struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	TxHash  string `json:"id"`
}

type Binaorder struct {
	Code          int    `json:"code"`
	Message       string `json:"msg"`
	symbol        string `json:"symbol"`
	OrderId       uint64 `json:"orderId"`
	ClientOrderId string `json:"clientOrderId"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	Status        string `json:"status"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	StopPrice     string `json:"stopPrice"`
	IcebergQty    string `json:"icebergQty"`
	Time          uint64 `json:"time"`
}
