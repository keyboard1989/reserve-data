package common

const TRUNC_LENGTH int = 600

func TruncStr(src []byte) []byte {
	if len(src) > TRUNC_LENGTH {
		return append(src[0:TRUNC_LENGTH], []byte("...")...)
	} else {
		return src
	}
}
