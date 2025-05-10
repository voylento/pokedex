package pokemonapi

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func bytesToExploreResponse(data []byte) (*ExploreResponse, error) {
	var exploreResponse ExploreResponse
	err := json.Unmarshal(data, &exploreResponse)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling explore response from PokeCache: %w", err)
	}
	return &exploreResponse, nil
}

func exploreResponseToBytes(exploreResponse *ExploreResponse) ([]byte, error) {
	bytes, err := json.Marshal(exploreResponse) 
	if err != nil {
		return nil, fmt.Errorf("error caching explore response: %w", err)
	}
	return bytes, nil
}

func GetExploreArea(url string) (*ExploreResponse, error) {
	if url == "" {
		return nil, fmt.Errorf("explore url not provided %d", 1)
	}

	if data, ok := PokeCache.Get(url); ok {
		return bytesToExploreResponse(data)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected HTTP Status code: %d", resp.StatusCode)
	}

	var exploreResponse ExploreResponse
	err = json.NewDecoder(resp.Body).Decode(&exploreResponse)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	bytes, err := exploreResponseToBytes(&exploreResponse)
	if err != nil {
	} else {
		PokeCache.Add(url, bytes)
	}
	
	return &exploreResponse, nil
}
