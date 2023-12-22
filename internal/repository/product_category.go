package repository

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

func (p ProductCategory) String() string {
	return [...]string{
		"",
		"Фрукты и ягоды",
		"Овощи и грибы",
		"Молоко, молочные продукты и их заменители",
		"Орехи, бобовые и семена",
		"Напитки",
		"Специи, травы, соусы и спреды",
		"Хлеб, злаки, макаронные изделия",
	}[p]
}
