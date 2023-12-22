package repository

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

func (c CarbType) String() string {
	return [...]string{
		"",
		"Фруктаны",
		"Фруктоза",
		"Сорбитол",
		"Маннитол",
		"Галактаны",
		"Галактоолигосахариды (ГОС)",
		"Лактоза",
		"Талактаны",
	}[c]
}
