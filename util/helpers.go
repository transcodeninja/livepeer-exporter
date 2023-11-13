// Package util provides various helper functions for the Livepeer exporter.
package util

import (
	"log"
	"math"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// BoolToFloat64 converts a bool to a float64.
// If the input bool is true, it returns 1.0; otherwise, it returns 0.0.
func BoolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// StringToFloat64 parses a string to a float64.
// If the string cannot be parsed, it returns an error.
func StringToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Error parsing value %v: %v", s, err)
		return 0.0, err
	}
	return f, nil
}

// ParseFloatAndSetGauge parses a string to a float64 and sets the value of the given gauge.
func ParseFloatAndSetGauge(value string, gauge prometheus.Gauge) {
	parsed, err := StringToFloat64(value)
	if err != nil {
		log.Printf("Error parsing value %v: %v", value, err)
		return
	}
	gauge.Set(parsed)
}

// Round rounds a float64 to a given number of decimal places.
func Round(value float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(value*shift) / shift
}
