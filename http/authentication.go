package http

type Authentication interface {
	KNSign(message string) string
	KNReadonlySign(message string) string
}
