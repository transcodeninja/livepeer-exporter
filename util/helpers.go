// Package util provides various helper functions for the Livepeer exporter.
package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

const graphQLEndpoint = "https://api.thegraph.com/subgraphs/name/livepeer/arbitrum-one"

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
func SetFloatFromStr(dest *float64, source string) {
	temp, err := StringToFloat64(source)
	if err != nil {
		log.Printf("Error parsing string to float: %v", err)
		return
	}
	*dest = temp
}

// getEnvVarDuration retrieves a duration from an environment variable.
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Fatalf("failed to parse '%s' environment variable: %v", key, err)
	}
	return value
}

// graphQLRequest represents the structure of the GraphQL API request used in IsOrchestrator.
type GraphQLRequest struct {
	Query string `json:"query"`
}

// graphqlResponse represents the structure of the GraphQL API response used in IsOrchestrator.
type graphQLResponse struct {
	Data struct {
		Transcoder struct {
			Typename string `json:"__typename"`
		}
	}
}

// sendGraphQLRequest sends a GraphQL request and returns the response body.
func sendGraphQLRequest(query string) ([]byte, error) {
	request := GraphQLRequest{
		Query: query,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(graphQLEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("received non-OK response from GraphQL endpoint")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseBody, nil
}

// IsOrchestrator checks if a given address is an Livepeer orchestrator.
func IsOrchestrator(id string) (bool, error) {
	query := fmt.Sprintf(`{
        transcoder(id: "%s") {
            __typename
        }
    }`, id)

	responseBody, err := sendGraphQLRequest(query)
	if err != nil {
		return false, err
	}

	var response graphQLResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data.Transcoder.Typename == "Transcoder", nil
}

// delegatorRequest represents the structure of the GraphQL API request used in IsDelegator.
type delegatorResponse struct {
	Data struct {
		Delegator struct {
			Typename string `json:"__typename"`
		}
	}
}

// IsDelegator checks if a given address is an Livepeer delegator.
func IsDelegator(id string) (bool, error) {
	query := fmt.Sprintf(`{
        delegator(id: "%s") {
            __typename
        }
    }`, id)

	responseBody, err := sendGraphQLRequest(query)
	if err != nil {
		return false, err
	}

	var response delegatorResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data.Delegator.Typename == "Delegator", nil
}
