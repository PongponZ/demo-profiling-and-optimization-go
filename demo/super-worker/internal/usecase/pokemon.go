package usecase

import (
	"math/rand"

	"github.com/PongponZ/demo-profiling-and-optimization-go/super-worker/internal/entity"
	"github.com/PongponZ/demo-profiling-and-optimization-go/super-worker/internal/repo"
)

type PokemonUsecase struct {
	repo *repo.PokemonRepo
}

func NewPokemonUsecase(repo *repo.PokemonRepo) *PokemonUsecase {
	return &PokemonUsecase{
		repo: repo,
	}
}

func (u *PokemonUsecase) GeneratePokemon(name string) entity.Pokemon {
	abilities := u.repo.FetchAbility()
	dna := u.GenerateDNA()
	stats := u.GenerateStats(dna)
	return entity.Pokemon{
		Name:      name,
		DNA:       dna,
		Abilities: abilities,
		Stats:     stats,
	}
}

func (u *PokemonUsecase) GenerateDNA() string {
	dna := ""
	base := "ATCG"

	for range 10000 {
		dna += string(base[rand.Intn(len(base))])
	}

	return dna
}

func (u *PokemonUsecase) GenerateStats(dna string) entity.Stats {
	base := len(dna)

	if base > 2 {
		base = base / 2
	}

	baseLv := rand.Intn(base) + 1
	hp := rand.Intn(base) + 1
	attack := rand.Intn(base) + 1
	defense := rand.Intn(base) + 1
	specialAttack := rand.Intn(base) + 1
	specialDefense := rand.Intn(base) + 1
	speed := rand.Intn(base) + 1

	return entity.Stats{
		BaseLv:         baseLv,
		HP:             hp,
		Attack:         attack,
		Defense:        defense,
		SpecialAttack:  specialAttack,
		SpecialDefense: specialDefense,
		Speed:          speed,
	}
}
