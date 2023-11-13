// Package util provides various helper functions for the Livepeer exporter.
package util

import (
	"log"
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

// ParseFloatAndSetGauge parses a string to a float64 and sets the value of the given gauge.
func ParseFloatAndSetGauge(value string, gauge prometheus.Gauge) {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing value %v: %v", value, err)
		return
	}
	gauge.Set(parsed)
}
