package listeners

// toUint safely converts an interface{} to uint
// This handles both direct uint values and float64 values (from JSON deserialization)
func toUint(v interface{}) (uint, bool) {
	switch val := v.(type) {
	case uint:
		return val, true
	case uint64:
		return uint(val), true
	case uint32:
		return uint(val), true
	case int:
		if val >= 0 {
			return uint(val), true
		}
		return 0, false
	case int64:
		if val >= 0 {
			return uint(val), true
		}
		return 0, false
	case int32:
		if val >= 0 {
			return uint(val), true
		}
		return 0, false
	case float64:
		if val >= 0 {
			return uint(val), true
		}
		return 0, false
	case float32:
		if val >= 0 {
			return uint(val), true
		}
		return 0, false
	default:
		return 0, false
	}
}

// toInt safely converts an interface{} to int
// This handles both direct int values and float64 values (from JSON deserialization)
func toInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int64:
		return int(val), true
	case int32:
		return int(val), true
	case uint:
		return int(val), true
	case uint64:
		return int(val), true
	case uint32:
		return int(val), true
	case float64:
		return int(val), true
	case float32:
		return int(val), true
	default:
		return 0, false
	}
}
