package utils

func IsContainValue(arr []interface{}, target interface{}) bool {
	for _, v := range arr {
		// Use type assertions to compare elements
		switch t := target.(type) {
		case string:
			if s, ok := v.(string); ok && s == t {
				return true
			}
		case int:
			if i, ok := v.(int); ok && i == t {
				return true
			}
		// Add more cases for other types as needed
		default:
			return false // Unsupported type
		}
	}
	return false
}
