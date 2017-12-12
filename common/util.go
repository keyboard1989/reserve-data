package common

const TRUNC_LENGTH int = 600

func TruncStr(src []byte) []byte {
	if len(src) > TRUNC_LENGTH {
		result := string(src[0:TRUNC_LENGTH]) + "..."
		return []byte(result)
	} else {
		return src
	}
}
