package utils

import (
	"strings"
)

func IsEmpty(i interface{}) bool {
	switch v := i.(type) {
	case string:
		if strings.TrimSpace(v) == "" || strings.TrimSpace(v) == `""` {
			return true
		}
	case int64:
		if v == 0 {
			return true
		}
	case int32:
		if v == 0 {
			return true
		}
	case float64:
		if v == 0 {
			return true
		}
	case []byte:
		if len(v) <= 0 {
			return true
		}
	default:
		if v == nil {
			return true
		}
		return false
	}
	return false
}
