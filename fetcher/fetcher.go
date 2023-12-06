// Package fetcher provides a utility for fetching and decoding JSON data from a remote server.
package fetcher

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Fetcher fetches JSON data from a specified URL and unmarshals it into a provided struct.
type Fetcher struct {
	URL  string      // URL to fetch data from.
	Data interface{} // Target struct to unmarshal data into.
}

// FetchData fetches JSON data from the Fetcher's URL and unmarshals it into the Fetcher's Data field.
// It returns an error if there was an issue fetching the data, if the HTTP status code is not 200,
// or if there was an issue decoding the response body.
func (f *Fetcher) FetchData() error {
	// Fetch data from the URL.
	resp, err := http.Get(f.URL)
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

// FetchDataWithBody fetches JSON data from the Fetcher's URL with the provided data and unmarshals
// it into the Fetcher's Data field. It returns an error if there was an issue fetching the data, if
// the HTTP status code is not 200, or if there was an issue decoding the response body.
func (f *Fetcher) FetchDataWithBody(data string) error {
	// Create a new request with the provided data.
	req, err := http.NewRequest("POST", f.URL, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")

	// Create a client with a timeout.
	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching data from '%s': %w", f.URL, err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Create a gzip reader.
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Decode the response body.
	dec := json.NewDecoder(gzipReader)
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

	resp, err := http.Post(f.URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error making GraphQL request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK response status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&f.Data)
	if err != nil {
		return fmt.Errorf("error decoding JSON response: %v", err)
	}

	return nil
}
