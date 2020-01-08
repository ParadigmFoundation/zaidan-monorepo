package utils

import "fmt"

func AnyToBytes(v interface{}) []byte {
	var bz []byte
	switch t := v.(type) {
	case []uint8:
		bz = t
	case string:
		bz = []byte(t)
	default:
		panic(fmt.Sprintf("can't convert %T to []byte", t))
	}
	return bz
}
