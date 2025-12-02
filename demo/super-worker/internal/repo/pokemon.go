package repo

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type PokemonRepo struct {
	url string
}

func NewPokemon(url string) *PokemonRepo {
	return &PokemonRepo{
		url: url,
	}
}

func (r *PokemonRepo) FetchAbility() map[string]int {
	ability := map[string]int{}

	for range rand.Intn(100) {
		response, err := http.Get(r.url)
		if err != nil {
			return nil
		}
		defer response.Body.Close()

		var abilities map[string]int
		err = json.NewDecoder(response.Body).Decode(&abilities)
		if err != nil {
			return nil
		}
		for k, v := range abilities {
			ability[k] += v
		}

		// simulate network latency
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	return ability
}
