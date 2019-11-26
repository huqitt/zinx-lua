package utils

import (
	"strconv"
)

func ToString(arg interface{}) string {
	result := ""
	switch val := arg.(type) {
	case string:
		result += val
	case bool:
		if val {
			result += "true"
		} else {
			result += "false"
		}
	case uint8:
		result += strconv.FormatUint(uint64(val), 10)
	case uint32:
		result += strconv.FormatUint(uint64(val), 10)
	case uint64:
		result += strconv.FormatUint(val, 10)
	case int32:
		result += strconv.Itoa(int(val))
	case int64:
		result += strconv.FormatInt(val, 10)
	case float32:
		result += strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		result += strconv.FormatFloat(float64(val), 'f', -1, 64)
	case []byte:
		result += "nil"
	default:
		result += strconv.Itoa(val.(int))
	}
	return result
}