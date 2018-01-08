package okex

type Signer interface {
	GetOkexKey() string
	OkexSign(msg string) string
}
