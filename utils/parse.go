package utils

import "strconv"

// ParseFloat safely converts string → float64
func ParseFloat(val string) float64 {
	if val == "" {
		return 0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return f
}
