package pokemonapi

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func bytesToLocationResponse(data []byte) (*LocationAreaResponse, error) {
	var locationResponse LocationAreaResponse
	err := json.Unmarshal(data, &locationResponse)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response from PokeCache: %w", err)
	}
	return &locationResponse, nil
}

func locationResponseToBytes(locationResponse *LocationAreaResponse) ([]byte, error) {
	bytes, err := json.Marshal(locationResponse) 
	if err != nil {
		return nil, fmt.Errorf("error caching response: %w", err)
	}
	return bytes, nil
}

func GetLocationAreas(url string) (*LocationAreaResponse, error) {
	if url == "" {
		return nil, fmt.Errorf("location-area url not provided %d", 1)
	}

	if data, ok := PokeCache.Get(url); ok {
		return bytesToLocationResponse(data)
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

	bytes, err := locationResponseToBytes(&locationResponse)
	if err != nil {
	} else {
		PokeCache.Add(url, bytes)
	}
	
	return &locationResponse, nil
}
