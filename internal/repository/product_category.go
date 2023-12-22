package repository

import "fmt"

type ProductCategory int

const (
	FruitsAndBerries ProductCategory = iota + 1
	VegetablesAndMushrooms
	MilkDairyProductsAndTheirSubstitutes
	NutsLegumesAndSeeds
	Drinks
	SpicesHerbsSaucesAndSpreads
	BreadCerealsPasta
)

var productCategoryMap = map[string]ProductCategory{
	"Фрукты и ягоды": FruitsAndBerries,
	"Овощи и грибы":  VegetablesAndMushrooms,
	"Молоко, молочные продукты и их заменители": MilkDairyProductsAndTheirSubstitutes,
	"Орехи, бобовые и семена":                   NutsLegumesAndSeeds,
	"Напитки": Drinks,
	"Специи, травы, соусы и спреды":    SpicesHerbsSaucesAndSpreads,
	"Хлеб, злаки и макаронные изделия": BreadCerealsPasta,
}

func StringToProductCategory(s string) (ProductCategory, error) {
	category, exists := productCategoryMap[s]
	if !exists {
		return 0, fmt.Errorf("unknown category: %s", s)
	}
	return category, nil
}
