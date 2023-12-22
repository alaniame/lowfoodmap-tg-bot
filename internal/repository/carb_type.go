package repository

import "fmt"

type CarbType int

const (
	Fructans CarbType = iota + 1
	Fructose
	Sorbitol
	Mannitol
	Galactans
	Galactooligosaccharides
	Lactose
	Talaktans
)

var carbTypeMap = map[string]CarbType{
	"Фруктаны":  Fructans,
	"Фруктоза":  Fructose,
	"Сорбитол":  Sorbitol,
	"Маннитол":  Mannitol,
	"Галактаны": Galactans,
	"Галактоолигосахариды (ГОС)": Galactooligosaccharides,
	"Лактоза":   Lactose,
	"Талактаны": Talaktans,
}

func StringToCarbType(s string) (CarbType, error) {
	if s == "" {
		return 0, nil
	}
	carbType, exists := carbTypeMap[s]
	if !exists {
		return 0, fmt.Errorf("unknown type: %s", s)
	}
	return carbType, nil
}
