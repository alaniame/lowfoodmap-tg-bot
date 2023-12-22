package repository

import (
	"fmt"
	"strings"
)

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
	"фруктаны":  Fructans,
	"фруктоза":  Fructose,
	"сорбитол":  Sorbitol,
	"маннитол":  Mannitol,
	"галактаны": Galactans,
	"галактоолигосахариды(гос)": Galactooligosaccharides,
	"лактоза":   Lactose,
	"талактаны": Talaktans,
}

func StringToCarbTypes(s string) ([]CarbType, error) {
	if s == "" {
		return nil, nil
	}
	var carbTypes []CarbType
	for _, carbStr := range strings.Split(s, " ") {
		normalizedCarbStr := strings.ToLower(carbStr)
		carbType, exists := carbTypeMap[normalizedCarbStr]
		if !exists {
			return nil, fmt.Errorf("unknown type: %s", carbStr)
		}
		carbTypes = append(carbTypes, carbType)
	}
	return carbTypes, nil
}
