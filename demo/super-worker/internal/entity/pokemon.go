package entity

type Pokemon struct {
	Name      string         `json:"name"`
	DNA       string         `json:"dna"`
	Ability   []string       `json:"ability"`
	Stats     Stats          `json:"stats"`
	Abilities map[string]int `json:"abilities"`
}

type Stats struct {
	BaseLv         int `json:"base_lv"`
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}
