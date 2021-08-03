package util

func IsNumberType(inter interface{}) bool {
	if inter == nil {
		return false
	}
	switch inter.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}
