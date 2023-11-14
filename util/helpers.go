// Package util provides various helper functions for the Livepeer exporter.
package util

import (
	"log"
	"math"
	"strconv"
)

// BoolToFloat64 converts a bool to a float64.
// If the input bool is true, it returns 1.0; otherwise, it returns 0.0.
func BoolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// Round rounds a float64 to a given number of decimal places.
func Round(value float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(value*shift) / shift
}

// StringToFloat64 parses a string to a float64.
// If the string cannot be parsed, it returns an error.
func StringToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Error parsing value %v: %v", s, err)
	}
	return f, err
}

// SetFloatFromStr sets the value of a float64 pointer from a string.
// If the string cannot be parsed to a float64, it logs an error and returns.
// If significance is not -1, it rounds the float64 to the given number of decimal places.
func SetFloatFromStr(dest *float64, source string, significance int) {
	temp, err := StringToFloat64(source)
	if err != nil {
		log.Printf("Error parsing string to float: %v", err)
		return
	}
	if significance != -1 {
		temp = Round(temp, significance)
	}
	*dest = temp
}
