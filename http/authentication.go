package http

type Authentication interface {
	KNSign(message string) string
}
