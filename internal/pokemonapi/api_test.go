package pokemonapi

import (
	"testing"
	"fmt"
)

func TestEmpty(t *testing.T) {
	cases := []struct {
		key string
		val []byte
	}{
	}

	for i, _ := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T){
		})
	}
}

