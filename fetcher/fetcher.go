// Package fetcher provides a utility for fetching and decoding JSON data from a remote server.
package fetcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Fetcher fetches JSON data from a specified URL and unmarshals it into a provided struct.
type Fetcher struct {
	URL     string      // URL to fetch data from.
	Data    interface{} // Target struct to unmarshal data into.
	Headers http.Header // Headers to send with the request.
}

// FetchData fetches JSON data from the Fetcher's URL and unmarshals it into the Fetcher's Data field.
// It returns an error if there was an issue fetching the data, if the HTTP status code is not 200,
// or if there was an issue decoding the response body.
func (f *Fetcher) FetchData() error {
	// Create a new request.
	req, err := http.NewRequest("GET", f.URL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Add additional headers, if any.
	if f.Headers != nil {
		for name, values := range f.Headers {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
	}

	// Create a client and send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching data from '%s': %w", f.URL, err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Decode the response body directly into the Fetcher's Data field.
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&f.Data); err != nil {
		return fmt.Errorf("error decoding response body from '%s': %w", f.URL, err)
	}

	return nil
}

// FetchGraphQLData fetches GraphQL data from the Fetcher's URL with the provided query and unmarshals
// it into the Fetcher's Data field. It returns an error if there was an issue fetching the data, if
// the HTTP status code is not 200, or if there was an issue decoding the response body.
func (f *Fetcher) FetchGraphQLData(query string) error {
	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		return fmt.Errorf("error creating request body: %v", err)
	}

	// Create a new request with the provided data.
	req, err := http.NewRequest("POST", f.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add additional headers, if any.
	if f.Headers != nil {
		for name, values := range f.Headers {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
	}

	// Send the request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making GraphQL request gtom '%s': %w", f.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&f.Data); err != nil {
		return fmt.Errorf("error decoding response body from '%s': %w", f.URL, err)
	}

	return nil
}
