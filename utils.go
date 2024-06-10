package main

import (
	"reflect"
	"strconv"
)

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	switch v := object.(type) {
	case bool:
		return v
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(object).Int() != 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(object).Uint() != 0
	case float32, float64:
		return reflect.ValueOf(object).Float() != 0
	case string:
		return v != ""
		// Default case, return true
	default:
		// For any other type, you can define your own truthy logic
		return true
	}
}

func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func castToFloat(a interface{}) float64 {
	if val, ok := a.(float64); ok {
		// If `right` is already a float64, simply negate it and return
		return val
	}
	if val, ok := a.(int); ok {
		return float64(val)
	}
	if val, ok := a.(string); ok {
		// Convert the string to float64 using strconv.ParseFloat()
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			// Handle error if conversion fails
			// For simplicity, let's return 0 in case of error
			return 0
		}
		return floatValue
		// Return the negation of the float64 value
	}
	return 0
}
