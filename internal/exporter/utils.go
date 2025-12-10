package exporter

import (
	"fmt"
)

func boolToFloat(v interface{}) float64 {
	switch t := v.(type) {
	case bool:
		if t {
			return 1
		}
		return 0
	case string:
		if t == "true" || t == "True" {
			return 1
		}
		return 0
	case float64:
		return t
	case int:
		return float64(t)
	default:
		return 0
	}
}

func numberToFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case string:
		// try parse
		return 0
	default:
		return 0
	}
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%v", t)
	case int:
		return fmt.Sprintf("%v", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func boolOrNumberToFloat(v interface{}) float64 {
	switch t := v.(type) {
	case bool:
		if t {
			return 1
		}
		return 0
	case float64:
		return t
	case int:
		return float64(t)
	case string:
		if t == "true" || t == "True" {
			return 1
		}
		// try parse numeric string
		return 0
	default:
		return 0
	}
}
