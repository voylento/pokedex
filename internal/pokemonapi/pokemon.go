package pokemonapi

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func bytesToPokemon(data []byte) (*Pokemon, error) {
	var pokemon Pokemon
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling pokemon response from PokeCache: %w", err)
	}
	return &pokemon, nil
}

func pokemonToBytes(pokemon *Pokemon) ([]byte, error) {
	bytes, err := json.Marshal(pokemon) 
	if err != nil {
		return nil, fmt.Errorf("error caching pokemon response: %w", err)
	}
	return bytes, nil
}

func GetPokemon(url string) (*Pokemon, error) {
	if url == "" {
		return nil, fmt.Errorf("pokemon url not provided %d", 1)
	}

	if data, ok := PokeCache.Get(url); ok {
		return bytesToPokemon(data)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected HTTP Status code: %d", resp.StatusCode)
	}

	var pokemon Pokemon

	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	bytes, err := pokemonToBytes(&pokemon)
	if err != nil {
	} else {
		PokeCache.Add(url, bytes)
	}
	
	return &pokemon, nil
}
