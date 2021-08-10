package utils

func Quot(a int64, b int64) int64 {
	return (a - (a % b)) / b
}

func GetFilledArray(length int, val byte) []byte {
	var res []byte
	for len(res) < length {
		res = append(res, val)
	}
	return res
}
