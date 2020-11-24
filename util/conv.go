package util

import (
	"strconv"
)

func ToString(v interface{}) string {
	switch vv := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(vv)
	case string:
		return vv
	case int:
		return strconv.Itoa(vv)
	case int64:
		return strconv.FormatInt(vv, 10)
	case bool:
		return strconv.FormatBool(vv)
	}
	return ""
}
