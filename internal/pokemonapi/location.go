package pokemonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationAreaResponse struct {
	Count 		int									`json:"count"`
	Next			interface{} 				`json:"next"`
	Previous	interface{}					`json:"previous"`
	Results		[]LocationAreaItem 	`json:"results"`
}

type LocationAreaItem struct {
	Name		string	`json:"name"`
	URL			string	`json:"url"`
}

func GetLocationAreas(url string) (*LocationAreaResponse, error) {
	if url == "" {
		return nil, fmt.Errorf("location-area url not provided %d", 1)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected HTTP Status code: %d", resp.StatusCode)
	}

	var locationResponse LocationAreaResponse
	err = json.NewDecoder(resp.Body).Decode(&locationResponse)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	return &locationResponse, nil
}
