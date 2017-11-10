package bittrex

type Signer interface {
	GetBittrexKey() string
	BittrexSign(msg string) string
}
