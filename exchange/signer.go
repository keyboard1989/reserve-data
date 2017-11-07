package exchange

type Signer interface {
	GetLiquiKey() string
	GetBittrexKey() string
	GetBinanceKey() string
	LiquiSign(msg string) string
	BittrexSign(msg string) string
	BinanceSign(msg string) string
}
