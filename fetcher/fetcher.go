// Package fetcher provides a utility for fetching and decoding JSON data from a remote server.
package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
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
